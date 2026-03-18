package use_else_if

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseElseIfRule detects else blocks containing only a nested if statement
// that can be simplified to else if.
type UseElseIfRule struct {
	skipBalanced bool
	configured   bool
}

const defaultSkipBalanced = true

// Configure validates and applies the rule configuration.
func (r *UseElseIfRule) Configure(arguments lint.Arguments) error {
	r.skipBalanced = defaultSkipBalanced

	if len(arguments) < 1 {
		r.configured = true
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "useElseIf" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "skipBalanced") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skipBalanced in "useElseIf" rule; need bool but got %T`, v)
			}
			r.skipBalanced = val
		}
	}

	r.configured = true
	return nil
}

// Apply applies the rule to given file.
func (r *UseElseIfRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.skipBalanced = defaultSkipBalanced
		r.configured = true
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseElseIf{
		onFailure:    onFailure,
		skipBalanced: r.skipBalanced,
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseElseIfRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.skipBalanced = defaultSkipBalanced
		r.configured = true
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseElseIf{
		onFailure:    onFailure,
		skipBalanced: r.skipBalanced,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseElseIfRule) Name() string {
	return "useElseIf"
}

// Group returns the rule group.
func (*UseElseIfRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseElseIfRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseElseIf struct {
	onFailure    func(lint.Failure)
	skipBalanced bool
}

func (w *lintUseElseIf) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Check if the else branch is a block containing only a single if statement
	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok {
		return w
	}

	if len(elseBlock.List) != 1 {
		return w
	}

	innerIf, ok := elseBlock.List[0].(*ast.IfStmt)
	if !ok {
		return w
	}

	// If skipBalanced is true, skip cases where both the outer if body
	// and the inner if body contain only a single statement each
	if w.skipBalanced {
		if len(ifStmt.Body.List) == 1 && len(innerIf.Body.List) == 1 {
			return w
		}
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryStyle,
		Node:       ifStmt.Else,
		Failure:    "can replace else block with else if",
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
