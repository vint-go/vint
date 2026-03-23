package no_off_by_one_error

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoOffByOneErrorRule detects off-by-one errors such as accessing s[len(s)]
// instead of s[len(s)-1].
type NoOffByOneErrorRule struct{}

// Apply applies the rule to given file.
func (r *NoOffByOneErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintOffByOne{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoOffByOneErrorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintOffByOne{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoOffByOneErrorRule) Name() string {
	return "noOffByOneError"
}

// Group returns the rule group.
func (*NoOffByOneErrorRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoOffByOneErrorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintOffByOne struct {
	onFailure func(lint.Failure)
}

func (w *lintOffByOne) Visit(node ast.Node) ast.Visitor {
	indexExpr, ok := node.(*ast.IndexExpr)
	if !ok {
		return w
	}

	// Check if the index is a call to len() with the same argument as the indexed expression.
	call, ok := indexExpr.Index.(*ast.CallExpr)
	if !ok {
		return w
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "len" {
		return w
	}

	if len(call.Args) != 1 {
		return w
	}

	// Compare the indexed expression with the len() argument.
	if exprEqual(indexExpr.X, call.Args[0]) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       indexExpr,
			Failure:    "off-by-one error: index equals length of the container",
		})
	}

	return w
}

// exprEqual checks if two AST expressions are structurally equal
// (for simple cases like identifiers and selector expressions).
func exprEqual(a, b ast.Expr) bool {
	switch x := a.(type) {
	case *ast.Ident:
		y, ok := b.(*ast.Ident)
		return ok && x.Name == y.Name
	case *ast.SelectorExpr:
		y, ok := b.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		return x.Sel.Name == y.Sel.Name && exprEqual(x.X, y.X)
	case *ast.IndexExpr:
		y, ok := b.(*ast.IndexExpr)
		if !ok {
			return false
		}
		return exprEqual(x.X, y.X) && exprEqual(x.Index, y.Index)
	case *ast.BasicLit:
		y, ok := b.(*ast.BasicLit)
		return ok && x.Kind == y.Kind && x.Value == y.Value
	}
	return false
}
