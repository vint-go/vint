package migrate

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/vint-go/vint/lint"
)

var nolintRegexp = regexp.MustCompile(`//nolint(?::(\S+))?`)

// nolintDirective represents a parsed //nolint comment in source code.
type nolintDirective struct {
	line        int      // 1-based line number
	linters     []string // linter names (empty = bare //nolint = all)
	colStart    int      // byte offset of //nolint within the line
	fullMatch   string   // the full matched text
	isInline    bool     // true if code precedes the //nolint on the same line
	explanation string   // text after // following the nolint (if any)
}

// nolintScope describes which lines a nolint directive covers.
type nolintScope struct {
	isFileLevel bool // true if the directive is before the package keyword
	startLine   int  // first line to check for failures (inclusive)
	endLine     int  // last line to check for failures (inclusive)
}

// NolintConvertResult holds the result of converting a single file.
type NolintConvertResult struct {
	OriginalPath string
	NewContent   string
	Directives   int // number of nolint directives converted
}

// ConvertNolintInFile reads a Go file, finds //nolint directives, runs relevant
// vint rules to determine which actually fire, and replaces nolint with vint-ignore.
// ruleConfigs contains the migrated configuration for each rule (keyed by full vint path).
func ConvertNolintInFile(filePath string, registry *RuleRegistry, ruleConfigs map[string]VintRuleConfig) (*NolintConvertResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return ConvertNolintInSource(filePath, string(content), registry, ruleConfigs)
}

// ConvertNolintInSource converts nolint directives in the given source string.
// ruleConfigs contains the migrated configuration for each rule (keyed by full vint path).
func ConvertNolintInSource(filePath string, contentStr string, registry *RuleRegistry, ruleConfigs map[string]VintRuleConfig) (*NolintConvertResult, error) {
	lines := strings.Split(contentStr, "\n")

	// Find all nolint directives.
	directives := findNolintDirectives(lines)
	if len(directives) == 0 {
		return &NolintConvertResult{
			OriginalPath: filePath,
			NewContent:   contentStr,
		}, nil
	}

	// Strip nolint comments from source so rules can fire.
	strippedLines := stripNolintComments(lines, directives)
	strippedSource := strings.Join(strippedLines, "\n")

	// Determine which rules to run.
	rulesToRun := collectRulesToRun(directives, registry)

	// Run the rules on the stripped source, using migrated configuration.
	var failures []lint.Failure
	var parsedFile *lint.File
	if len(rulesToRun) > 0 {
		failures, parsedFile = runRulesOnSource(filePath, []byte(strippedSource), rulesToRun, ruleConfigs)
	}

	// Build failure lookup: line → set of rule names.
	failuresByLine := make(map[int]map[string]bool)
	for _, f := range failures {
		line := f.Position.Start.Line
		if failuresByLine[line] == nil {
			failuresByLine[line] = make(map[string]bool)
		}
		failuresByLine[line][f.RuleName] = true
	}

	// Compute nolint scopes using AST information.
	scopes := computeDirectiveScopes(directives, parsedFile, len(lines))

	// Build index of which lines have directives.
	directiveByLine := make(map[int]*nolintDirective)
	for i := range directives {
		directiveByLine[directives[i].line] = &directives[i]
	}

	// Build output top-down.
	var resultLines []string
	for lineNum := 1; lineNum <= len(lines); lineNum++ {
		origLine := lines[lineNum-1]
		d, hasDirective := directiveByLine[lineNum]

		if !hasDirective {
			resultLines = append(resultLines, origLine)
			continue
		}

		scope := scopes[d.line]

		// Determine which rules to suppress.
		candidateRules := matchRulesToDirective(d, failuresByLine, registry, scope)

		// Determine explanation text.
		explanation := d.explanation
		if explanation == "" {
			if len(d.linters) > 0 {
				explanation = fmt.Sprintf("migrated from nolint:%s", strings.Join(d.linters, ","))
			} else {
				explanation = "migrated from nolint"
			}
		}

		// Get indentation of original line.
		indent := leadingWhitespace(origLine)

		// Choose the correct vint-ignore variant.
		ignoreDirective := "vint-ignore"
		if scope.isFileLevel {
			ignoreDirective = "vint-ignore-all"
		}

		// Generate vint-ignore lines and insert before the code line.
		if len(candidateRules) == 0 {
			// No rules fired — the nolint was unnecessary.
			resultLines = append(resultLines,
				indent+fmt.Sprintf("// NOTE: nolint removed — no vint rules fired (was: %s)", d.fullMatch))
		} else {
			for _, rule := range candidateRules {
				resultLines = append(resultLines,
					indent+fmt.Sprintf("// %s %s: %s", ignoreDirective, rule, explanation))
			}
		}

		// Emit the code line with nolint stripped.
		if d.isInline {
			// Strip //nolint from the line, keep the code.
			resultLines = append(resultLines, strings.TrimRight(origLine[:d.colStart], " \t"))
		}
		// If not inline (the entire line was //nolint), skip emitting it —
		// the vint-ignore lines above replace it.
	}

	return &NolintConvertResult{
		OriginalPath: filePath,
		NewContent:   strings.Join(resultLines, "\n"),
		Directives:   len(directives),
	}, nil
}

// matchRulesToDirective finds which rules from the directive's linters actually fire
// within the scope of the directive. For aggregating rules (which require cross-file
// analysis and cannot fire from single-file Apply), we trust the mapping unconditionally.
func matchRulesToDirective(
	d *nolintDirective,
	failuresByLine map[int]map[string]bool,
	registry *RuleRegistry,
	scope nolintScope,
) []string {
	// ruleFiresInScope checks whether the given rule has a failure within the scope.
	ruleFiresInScope := func(ruleName string) bool {
		for line := scope.startLine; line <= scope.endLine; line++ {
			if failuresByLine[line] != nil && failuresByLine[line][ruleName] {
				return true
			}
		}
		return false
	}

	var candidateRules []string

	if len(d.linters) == 0 {
		// Bare //nolint — return all rules that fire within scope.
		seen := make(map[string]bool)
		for line := scope.startLine; line <= scope.endLine; line++ {
			for r := range failuresByLine[line] {
				if !seen[r] {
					seen[r] = true
					candidateRules = append(candidateRules, r)
				}
			}
		}
	} else {
		// Specific linters — only include rules mapped to those linters that fire.
		for _, linter := range d.linters {
			for _, mapped := range registry.RulesForLinter(linter) {
				if mapped.FullVintPath != "" {
					if isAggregatingRule(mapped.Rule) {
						candidateRules = append(candidateRules, mapped.FullVintPath)
					} else if ruleFiresInScope(mapped.FullVintPath) {
						candidateRules = append(candidateRules, mapped.FullVintPath)
					}
				}
			}
		}
	}

	sort.Strings(candidateRules)
	return candidateRules
}

// isAggregatingRule checks whether a rule implements lint.AggregatingRule.
func isAggregatingRule(r lint.Rule) bool {
	if r == nil {
		return false
	}
	_, ok := r.(lint.AggregatingRule)
	return ok
}

// findNolintDirectives scans source lines for //nolint comments.
func findNolintDirectives(lines []string) []nolintDirective {
	var directives []nolintDirective
	for i, line := range lines {
		loc := nolintRegexp.FindStringIndex(line)
		if loc == nil {
			continue
		}

		match := nolintRegexp.FindStringSubmatch(line)
		var linters []string
		if len(match) > 1 && match[1] != "" {
			linters = strings.Split(match[1], ",")
		}

		// Check if there's code before the //nolint (inline).
		prefix := strings.TrimSpace(line[:loc[0]])
		isInline := prefix != ""

		// Check for explanation after nolint directive.
		rest := line[loc[1]:]
		var explanation string
		if idx := strings.Index(rest, "//"); idx != -1 {
			explanation = strings.TrimSpace(rest[idx+2:])
		}

		directives = append(directives, nolintDirective{
			line:        i + 1, // 1-based
			linters:     linters,
			colStart:    loc[0],
			fullMatch:   match[0],
			isInline:    isInline,
			explanation: explanation,
		})
	}
	return directives
}

// stripNolintComments removes //nolint... from lines so rules can fire.
func stripNolintComments(lines []string, directives []nolintDirective) []string {
	stripped := make([]string, len(lines))
	copy(stripped, lines)
	for _, d := range directives {
		idx := d.line - 1
		loc := nolintRegexp.FindStringIndex(stripped[idx])
		if loc != nil {
			stripped[idx] = strings.TrimRight(stripped[idx][:loc[0]], " \t")
		}
	}
	return stripped
}

// collectRulesToRun determines which lint.Rule instances need to be run.
func collectRulesToRun(directives []nolintDirective, registry *RuleRegistry) []lint.Rule {
	seen := make(map[string]bool)
	var rules []lint.Rule

	for _, d := range directives {
		if len(d.linters) == 0 {
			// Bare //nolint — need all rules.
			return registry.AllRuleInstances()
		}
	}

	for _, d := range directives {
		for _, linter := range d.linters {
			for _, r := range registry.RuleInstancesForLinter(linter) {
				fullName := lint.FullRuleName(r)
				if !seen[fullName] {
					seen[fullName] = true
					rules = append(rules, r)
				}
			}
		}
	}
	return rules
}

// runRulesOnSource creates a lint.File from raw source and runs rules against it.
// ruleConfigs provides the migrated configuration for each rule so they fire correctly.
// Returns the failures and the parsed lint.File (for AST access); file may be nil on parse error.
func runRulesOnSource(filePath string, source []byte, rules []lint.Rule, ruleConfigs map[string]VintRuleConfig) ([]lint.Failure, *lint.File) {
	fset := token.NewFileSet()
	goVer := goversion.Must(goversion.NewVersion("1.22"))
	imp := lint.NewSharedImporter()
	pkg := lint.NewPackage(fset, goVer, imp)

	file, err := pkg.AddFile(filePath, source)
	if err != nil {
		return nil, nil
	}

	// Configure rules with migrated settings before running them.
	// Always call Configure() for configurable rules, even with empty options,
	// so they can perform initialization (e.g., building internal data structures).
	for _, rule := range rules {
		if cr, ok := rule.(lint.ConfigurableRule); ok {
			fullName := lint.FullRuleName(rule)
			if cfg, ok := ruleConfigs[fullName]; ok && len(cfg.Options) > 0 {
				_ = cr.Configure(lint.Arguments{cfg.Options})
			} else {
				_ = cr.Configure(nil)
			}
		}
	}

	var allFailures []lint.Failure
	for _, rule := range rules {
		fullName := lint.FullRuleName(rule)
		// Build arguments from migrated config options.
		var args lint.Arguments
		if cfg, ok := ruleConfigs[fullName]; ok && len(cfg.Options) > 0 {
			args = lint.Arguments{cfg.Options}
		}
		failures := rule.Apply(file, args)
		for i, f := range failures {
			if f.RuleName == "" {
				f.RuleName = fullName
			}
			if f.Node != nil {
				f.Position = lint.ToFailurePosition(f.Node.Pos(), f.Node.End(), file)
			}
			failures[i] = f
		}
		allFailures = append(allFailures, failures...)
	}

	return allFailures, file
}

// computeDirectiveScopes determines the effective line range for each nolint directive.
// For inline directives, the scope is just the directive's own line.
// For standalone directives before the package keyword, the scope is the entire file (file-level).
// For standalone directives above a statement, the scope covers the next AST node's line range.
func computeDirectiveScopes(directives []nolintDirective, file *lint.File, totalLines int) map[int]nolintScope {
	scopes := make(map[int]nolintScope, len(directives))

	// Determine the package declaration line.
	packageLine := math.MaxInt32
	if file != nil && file.AST != nil && file.AST.Package.IsValid() {
		packageLine = file.ToPosition(file.AST.Package).Line
	}

	// Collect sorted node line ranges from the AST for standalone scope resolution.
	var nodeRanges []nodeRange
	if file != nil && file.AST != nil {
		nodeRanges = collectNodeRanges(file)
	}

	for _, d := range directives {
		if d.isInline {
			// Inline: only suppress on this exact line.
			scopes[d.line] = nolintScope{startLine: d.line, endLine: d.line}
			continue
		}

		if d.line < packageLine {
			// File-level: before the package keyword → entire file.
			scopes[d.line] = nolintScope{
				isFileLevel: true,
				startLine:   1,
				endLine:     totalLines,
			}
			continue
		}

		// Standalone comment above a statement: find the next AST node.
		scope := nolintScope{startLine: d.line, endLine: d.line}
		for _, nr := range nodeRanges {
			if nr.start > d.line {
				scope.startLine = nr.start
				scope.endLine = nr.end
				break
			}
		}
		scopes[d.line] = scope
	}

	return scopes
}

// nodeRange represents the start and end line of an AST node.
type nodeRange struct {
	start int
	end   int
}

// collectNodeRanges walks the AST and collects line ranges for all
// declarations, specs, and statements, sorted by start line.
func collectNodeRanges(file *lint.File) []nodeRange {
	var ranges []nodeRange
	ast.Inspect(file.AST, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch n.(type) {
		case *ast.GenDecl, *ast.FuncDecl,
			*ast.ValueSpec, *ast.TypeSpec, // specs inside var()/const()/type() blocks
			*ast.AssignStmt, *ast.ExprStmt, *ast.ReturnStmt,
			*ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.TypeSwitchStmt,
			*ast.SelectStmt, *ast.GoStmt, *ast.DeferStmt,
			*ast.SendStmt, *ast.IncDecStmt, *ast.BranchStmt,
			*ast.LabeledStmt, *ast.DeclStmt:
			start := file.ToPosition(n.Pos()).Line
			end := file.ToPosition(n.End()).Line
			ranges = append(ranges, nodeRange{start: start, end: end})
		}
		return true
	})
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})
	return ranges
}

func leadingWhitespace(s string) string {
	for i, ch := range s {
		if ch != ' ' && ch != '\t' {
			return s[:i]
		}
	}
	return s
}
