package use_early_continue

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseEarlyContinueRule finds where nesting level could be reduced in loop bodies
// by inverting the condition and using continue.
type UseEarlyContinueRule struct {
	bodyWidth int
}

const defaultBodyWidth = 5

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *UseEarlyContinueRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.bodyWidth = defaultBodyWidth
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		width, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "useEarlyContinue" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.bodyWidth = int(width)
		return nil
	}

	r.bodyWidth = defaultBodyWidth
	for k, v := range argKV {
		if isRuleOption(k, "body-width") {
			width, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for body-width in "useEarlyContinue" rule; need int64 but got %T`, v)
			}
			r.bodyWidth = int(width)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *UseEarlyContinueRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.bodyWidth <= 0 {
		r.bodyWidth = defaultBodyWidth
	}

	var failures []lint.Failure
	w := &lintEarlyContinue{
		bodyWidth: r.bodyWidth,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseEarlyContinueRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if r.bodyWidth <= 0 {
		r.bodyWidth = defaultBodyWidth
	}

	var failures []lint.Failure
	w := &lintEarlyContinue{
		bodyWidth: r.bodyWidth,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseEarlyContinueRule) Name() string {
	return "useEarlyContinue"
}

// Group returns the rule group.
func (*UseEarlyContinueRule) Group() string {
	return "complexity"
}

// CacheTier returns the cache tier for this rule.
func (*UseEarlyContinueRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintEarlyContinue struct {
	bodyWidth int
	onFailure func(lint.Failure)
}

func (w *lintEarlyContinue) Visit(node ast.Node) ast.Visitor {
	// We look for for/range statements whose body contains exactly one statement
	// that is an if-statement with no else clause, and the if body has at least
	// bodyWidth statements.
	var body *ast.BlockStmt
	switch n := node.(type) {
	case *ast.ForStmt:
		body = n.Body
	case *ast.RangeStmt:
		body = n.Body
	default:
		return w
	}

	if body == nil || len(body.List) != 1 {
		return w
	}

	ifStmt, ok := body.List[0].(*ast.IfStmt)
	if !ok {
		return w
	}

	// Must have no else clause
	if ifStmt.Else != nil {
		return w
	}

	// Must have no init statement in the if
	if ifStmt.Init != nil {
		return w
	}

	// The if body must have at least bodyWidth statements
	if ifStmt.Body == nil || len(ifStmt.Body.List) < w.bodyWidth {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryComplexity,
		Failure:    "invert if condition and use continue to reduce nesting",
		Node:       ifStmt,
	})

	return w
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
