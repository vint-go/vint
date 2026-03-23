package migrate

import (
	"fmt"
	"go/token"
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
	if len(rulesToRun) > 0 {
		failures = runRulesOnSource(filePath, []byte(strippedSource), rulesToRun, ruleConfigs)
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

		// Determine which rules to suppress.
		candidateRules := matchRulesToDirective(d, failuresByLine, registry)

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

		// Generate vint-ignore lines and insert before the code line.
		if len(candidateRules) == 0 {
			// No rules fired — the nolint was unnecessary.
			resultLines = append(resultLines,
				indent+fmt.Sprintf("// NOTE: nolint removed — no vint rules fired (was: %s)", d.fullMatch))
		} else {
			for _, rule := range candidateRules {
				resultLines = append(resultLines,
					indent+fmt.Sprintf("// vint-ignore %s: %s", rule, explanation))
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

// matchRulesToDirective finds which rules from the directive's linters actually fire.
// For aggregating rules (which require cross-file analysis and cannot fire from
// single-file Apply), we trust the linter-to-rule mapping unconditionally.
func matchRulesToDirective(
	d *nolintDirective,
	failuresByLine map[int]map[string]bool,
	registry *RuleRegistry,
) []string {
	var candidateRules []string

	if len(d.linters) == 0 {
		// Bare //nolint — return all rules that fire on this line.
		if rules, ok := failuresByLine[d.line]; ok {
			for r := range rules {
				candidateRules = append(candidateRules, r)
			}
		}
	} else {
		// Specific linters — only include rules mapped to those linters that fire.
		for _, linter := range d.linters {
			for _, mapped := range registry.RulesForLinter(linter) {
				if mapped.FullVintPath != "" {
					// Aggregating rules (e.g. noDuplicateCode) use Collect/Finalize
					// and always return nil from Apply, so they cannot fire during
					// single-file analysis. Trust the mapping unconditionally.
					if isAggregatingRule(mapped.Rule) {
						candidateRules = append(candidateRules, mapped.FullVintPath)
					} else if failuresByLine[d.line] != nil && failuresByLine[d.line][mapped.FullVintPath] {
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
func runRulesOnSource(filePath string, source []byte, rules []lint.Rule, ruleConfigs map[string]VintRuleConfig) []lint.Failure {
	fset := token.NewFileSet()
	goVer := goversion.Must(goversion.NewVersion("1.22"))
	imp := lint.NewSharedImporter()
	pkg := lint.NewPackage(fset, goVer, imp)

	file, err := pkg.AddFile(filePath, source)
	if err != nil {
		return nil
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

	return allFailures
}

func leadingWhitespace(s string) string {
	for i, ch := range s {
		if ch != ' ' && ch != '\t' {
			return s[:i]
		}
	}
	return s
}
