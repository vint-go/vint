package lint

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"math"
	"regexp"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
)

// File abstraction used for representing files.
type File struct {
	Name    string
	Pkg     *Package
	content []byte
	AST     *ast.File
}

// IsTest returns if the file contains tests.
func (f *File) IsTest() bool { return strings.HasSuffix(f.Name, "_test.go") }

// IsImportable returns if the symbols defined in this file can be imported in other packages.
//
// Symbols from the package `main` or test files are not exported, so they cannot be imported.
func (f *File) IsImportable() bool {
	if f.IsTest() {
		// Test files cannot be imported.
		return false
	}

	if f.Pkg.IsMain() {
		// The package `main` cannot be imported.
		return false
	}

	return true
}

// Content returns the file's content.
func (f *File) Content() []byte {
	return f.content
}

// NewFile creates a new file.
func NewFile(name string, content []byte, pkg *Package) (*File, error) {
	f, err := parser.ParseFile(pkg.fset, name, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &File{
		Name:    name,
		content: content,
		Pkg:     pkg,
		AST:     f,
	}, nil
}

// ToPosition returns line and column for given position.
func (f *File) ToPosition(pos token.Pos) token.Position {
	return f.Pkg.fset.Position(pos)
}

// Render renders a node.
func (f *File) Render(x any) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, f.Pkg.fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

// CommentMap builds a comment map for the file.
func (f *File) CommentMap() ast.CommentMap {
	return ast.NewCommentMap(f.Pkg.fset, f.AST, f.AST.Comments)
}

var basicTypeKinds = map[types.BasicKind]string{
	types.UntypedBool:    "bool",
	types.UntypedInt:     "int",
	types.UntypedRune:    "rune",
	types.UntypedFloat:   "float64",
	types.UntypedComplex: "complex128",
	types.UntypedString:  "string",
}

// IsUntypedConst reports whether expr is an untyped constant,
// and indicates what its default type is.
// Scope may be nil.
func (f *File) IsUntypedConst(expr ast.Expr) (defType string, ok bool) {
	// Re-evaluate expr outside its context to see if it's untyped.
	// (An expr evaluated within, for example, an assignment context will get the type of the LHS.)
	exprStr := f.Render(expr)
	tv, err := types.Eval(f.Pkg.fset, f.Pkg.TypesPkg(), expr.Pos(), exprStr)
	if err != nil {
		return "", false
	}
	if b, ok := tv.Type.(*types.Basic); ok {
		if dt, ok := basicTypeKinds[b.Kind()]; ok {
			return dt, true
		}
	}

	return "", false
}

func (f *File) isMain() bool {
	return f.AST.Name.Name == "main"
}

const directiveSpecifyDisableReason = "specify-disable-reason"

type togetherApplier struct {
	rules       []WalkingRule
	file        *File
	args        Arguments
	allFailures []Failure
}

const walkingOptimizationOn = false

func (f *File) lint(rules []Rule, config Config, failures chan Failure, rc *rulecache.RuleCache, preCachedHits map[string]bool) error {
	rulesConfig := config.Rules
	_, mustSpecifyDisableReason := config.Directives[directiveSpecifyDisableReason]
	disabledIntervals := f.disabledIntervals(rules, mustSpecifyDisableReason, failures)

	// Collect sibling file contents for package-aware cache tiers (computed lazily).
	var siblingFiles map[string][]byte
	// Pre-computed sibling digest — computed once on first need, reused for
	// every subsequent package-aware rule on this file.
	var siblingDigest *[32]byte

	var walkingRules []WalkingRule = make([]WalkingRule, 0, len(rules))

	for _, currentRule := range rules {
		fullName := FullRuleName(currentRule)
		ruleConfig := rulesConfig[fullName]
		if ruleConfig.MustExclude(f.Name) {
			continue
		}

		// Skip rules already served from pre-parse cache check.
		if preCachedHits[fullName] {
			continue
		}

		// Determine cache tier and cacheability.
		tier := rulecache.TierCrossPackage // safe default
		cacheable := rc != nil
		if tr, ok := currentRule.(TieredRule); ok {
			tier = tr.CacheTier()
		}
		if ur, ok := currentRule.(UncacheableRule); ok && ur.Uncacheable() {
			cacheable = false
		}

		if walkingOptimizationOn {
			if wr, ok := currentRule.(WalkingRule); ok {
				walkingRules = append(walkingRules, wr)
				continue
			}
		}

		// Lazily collect sibling contents and compute the sibling digest
		// once for all package-aware (and above) rules on this file.
		if cacheable && tier >= rulecache.TierPackageAware {
			if siblingFiles == nil {
				siblingFiles = f.collectSiblingContents()
			}
			if siblingDigest == nil {
				d := rc.SiblingDigest(f.Name, siblingFiles)
				siblingDigest = &d
			}
		}

		// Try cache.
		if cacheable {
			if cached, hit := rc.Get(fullName, f.Name, tier, siblingFiles, siblingDigest); hit {
				for _, cf := range cached {
					failure := fromCachedFailure(cf)
					if failure.Confidence >= config.Confidence {
						failures <- failure
					}
				}
				continue
			}
		}

		// Cache miss — run the rule.
		currentFailures := currentRule.Apply(f, ruleConfig.Arguments)
		for idx, failure := range currentFailures {
			if failure.IsInternal() {
				return errors.New(failure.Failure)
			}

			if failure.RuleName == "" {
				failure.RuleName = fullName
			}
			if failure.Node != nil {
				failure.Position = ToFailurePosition(failure.Node.Pos(), failure.Node.End(), f)
			}
			currentFailures[idx] = failure
		}

		// Store in cache (post-filter).
		if cacheable {
			cached := make([]rulecache.CachedFailure, len(currentFailures))
			for i, fail := range currentFailures {
				cached[i] = toCachedFailure(fail)
			}
			rc.Put(fullName, f.Name, tier, siblingFiles, siblingDigest, cached)
		}

		currentFailures = f.filterFailures(currentFailures, disabledIntervals)
		for _, failure := range currentFailures {
			if failure.Confidence >= config.Confidence {
				failures <- failure
			}
		}
	}

	if len(walkingRules) > 0 {
		applier := &togetherApplier{
			rules: walkingRules,
			file:  f,
			args:  Arguments{},
		}
		ast.Walk(applier, f.AST)
		for _, failure := range applier.allFailures {
			if failure.IsInternal() {
				return errors.New(failure.Failure)
			}
		}
		applier.allFailures = f.filterFailures(applier.allFailures, disabledIntervals)
		for _, failure := range applier.allFailures {
			if failure.Confidence >= config.Confidence {
				failures <- failure
			}
		}
	}

	return nil
}

func (v *togetherApplier) Visit(node ast.Node) ast.Visitor {
	for _, rule := range v.rules {
		newFailures := rule.ApplyToNode(v.file, node, v.args)
		for idx, failure := range newFailures {
			if failure.RuleName == "" {
				failure.RuleName = FullRuleName(rule)
			}
			if failure.Node != nil {
				failure.Position = ToFailurePosition(failure.Node.Pos(), failure.Node.End(), v.file)
			}
			newFailures[idx] = failure
		}
		v.allFailures = append(v.allFailures, newFailures...)
	}
	return v
}

// collectSiblingContents returns the content of all files in the same package.
func (f *File) collectSiblingContents() map[string][]byte {
	files := f.Pkg.Files()
	result := make(map[string][]byte, len(files))
	for name, file := range files {
		result[name] = file.content
	}
	return result
}

// toCachedFailure converts a Failure to a CachedFailure for caching.
func toCachedFailure(f Failure) rulecache.CachedFailure {
	return rulecache.CachedFailure{
		Message:         f.Failure,
		RuleName:        f.RuleName,
		Category:        string(f.Category),
		StartFilename:   f.Position.Start.Filename,
		StartLine:       f.Position.Start.Line,
		StartColumn:     f.Position.Start.Column,
		EndFilename:     f.Position.End.Filename,
		EndLine:         f.Position.End.Line,
		EndColumn:       f.Position.End.Column,
		Confidence:      f.Confidence,
		ReplacementLine: f.ReplacementLine,
	}
}

// fromCachedFailure converts a CachedFailure back to a Failure.
func fromCachedFailure(cf rulecache.CachedFailure) Failure {
	return Failure{
		Failure:  cf.Message,
		RuleName: cf.RuleName,
		Category: FailureCategory(cf.Category),
		Position: FailurePosition{
			Start: token.Position{
				Filename: cf.StartFilename,
				Line:     cf.StartLine,
				Column:   cf.StartColumn,
			},
			End: token.Position{
				Filename: cf.EndFilename,
				Line:     cf.EndLine,
				Column:   cf.EndColumn,
			},
		},
		Confidence:      cf.Confidence,
		ReplacementLine: cf.ReplacementLine,
	}
}

type enableDisableConfig struct {
	enabled  bool
	position int
}

type disabledIntervalsMap = map[string][]DisabledInterval

const (
	directivePos = 1
	modifierPos  = 2
	rulesPos     = 3
	reasonPos    = 4
)

var directiveRegexp = regexp.MustCompile(`^//[\s]*revive:(enable|disable)(?:-(line|next-line))?(?::([^\s]+))?[\s]*(?: (.+))?$`)

func (f *File) disabledIntervals(rules []Rule, mustSpecifyDisableReason bool, failures chan Failure) disabledIntervalsMap {
	enabledDisabledRulesMap := map[string][]enableDisableConfig{}

	getEnabledDisabledIntervals := func() disabledIntervalsMap {
		result := disabledIntervalsMap{}

		for ruleName, disabledArr := range enabledDisabledRulesMap {
			ruleResult := []DisabledInterval{}
			for i := range disabledArr {
				interval := DisabledInterval{
					RuleName: ruleName,
					From: token.Position{
						Filename: f.Name,
						Line:     disabledArr[i].position,
					},
					To: token.Position{
						Filename: f.Name,
						Line:     math.MaxInt32,
					},
				}
				if i%2 == 0 {
					ruleResult = append(ruleResult, interval)
				} else {
					ruleResult[len(ruleResult)-1].To.Line = disabledArr[i].position
				}
			}
			result[ruleName] = ruleResult
		}

		return result
	}

	handleConfig := func(isEnabled bool, line int, name string) {
		existing, ok := enabledDisabledRulesMap[name]
		if !ok {
			existing = []enableDisableConfig{}
			enabledDisabledRulesMap[name] = existing
		}
		if (len(existing) > 1 && existing[len(existing)-1].enabled == isEnabled) ||
			(len(existing) == 0 && isEnabled) {
			return
		}
		existing = append(existing, enableDisableConfig{
			enabled:  isEnabled,
			position: line,
		})
		enabledDisabledRulesMap[name] = existing
	}

	handleRules := func(modifier string, isEnabled bool, line int, ruleNames []string) {
		for _, name := range ruleNames {
			switch modifier {
			case "line":
				handleConfig(isEnabled, line, name)
				handleConfig(!isEnabled, line, name)
			case "next-line":
				handleConfig(isEnabled, line+1, name)
				handleConfig(!isEnabled, line+1, name)
			default:
				handleConfig(isEnabled, line, name)
			}
		}
	}

	handleComment := func(c *ast.CommentGroup, line int) {
		comments := c.List
		for _, c := range comments {
			match := directiveRegexp.FindStringSubmatch(c.Text)
			if len(match) == 0 {
				continue
			}
			ruleNames := []string{}

			for name := range strings.SplitSeq(match[rulesPos], ",") {
				name = strings.Trim(name, "\n")
				if name != "" {
					ruleNames = append(ruleNames, name)
				}
			}

			mustCheckDisablingReason := mustSpecifyDisableReason && match[directivePos] == "disable"
			if mustCheckDisablingReason && strings.Trim(match[reasonPos], " ") == "" {
				failures <- Failure{
					Confidence: 1,
					RuleName:   directiveSpecifyDisableReason,
					Failure:    "reason of lint disabling not found",
					Position:   ToFailurePosition(c.Pos(), c.End(), f),
					Node:       c,
				}
				continue // skip this linter disabling directive
			}

			// TODO: optimize
			if len(ruleNames) == 0 {
				for _, rule := range rules {
					ruleNames = append(ruleNames, FullRuleName(rule))
				}
			}

			handleRules(match[modifierPos], match[directivePos] == "enable", line, ruleNames)
		}
	}

	for _, c := range f.AST.Comments {
		handleComment(c, f.ToPosition(c.End()).Line)
	}

	return getEnabledDisabledIntervals()
}

func (*File) filterFailures(failures []Failure, disabledIntervals disabledIntervalsMap) []Failure {
	result := []Failure{}
	for _, failure := range failures {
		fStart := failure.Position.Start.Line
		fEnd := failure.Position.End.Line
		intervals, ok := disabledIntervals[failure.RuleName]
		if !ok {
			result = append(result, failure)
			continue
		}

		include := true
		for _, interval := range intervals {
			intStart := interval.From.Line
			intEnd := interval.To.Line
			if (fStart >= intStart && fStart <= intEnd) ||
				(fEnd >= intStart && fEnd <= intEnd) {
				include = false
				break
			}
		}
		if include {
			result = append(result, failure)
		}
	}
	return result
}
