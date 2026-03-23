package no_discarded_append

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDiscardedAppendRule detects calls to append whose return value is discarded.
type NoDiscardedAppendRule struct{}

// Apply applies the rule to given file.
func (r *NoDiscardedAppendRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDiscardedAppend{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDiscardedAppendRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDiscardedAppend{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDiscardedAppendRule) Name() string {
	return "noDiscardedAppend"
}

// Group returns the rule group.
func (*NoDiscardedAppendRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDiscardedAppendRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDiscardedAppend struct {
	onFailure func(lint.Failure)
}

func (w *lintDiscardedAppend) Visit(node ast.Node) ast.Visitor {
	exprStmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return w
	}

	if containsAppendCall(exprStmt.X) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       exprStmt,
			Category:   lint.FailureCategoryLogic,
			Failure:    "result of append is discarded and has no effect",
		})
	}

	return w
}

// containsAppendCall checks if the expression is a call to the builtin append,
// or contains a nested call to append (e.g., append(append(s, 1), 2)).
func containsAppendCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name == "append" {
		return true
	}

	return false
}
