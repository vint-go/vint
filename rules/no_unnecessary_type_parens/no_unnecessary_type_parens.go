package no_unnecessary_type_parens

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryTypeParensRule detects unneeded parentheses inside type expressions
// and suggests removing them.
type NoUnnecessaryTypeParensRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryTypeParensRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintTypeParens{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryTypeParensRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintTypeParens{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryTypeParensRule) Name() string {
	return "noUnnecessaryTypeParens"
}

// Group returns the rule group.
func (*NoUnnecessaryTypeParensRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryTypeParensRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTypeParens struct {
	onFailure func(lint.Failure)
}

func (w *lintTypeParens) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.StarExpr:
		// Pointer type: *(<type>) -> check if X is parenthesized
		w.checkParenExpr(n.X)
	case *ast.ArrayType:
		// Array/slice element type: [N](<type>) or [](<type>)
		w.checkParenExpr(n.Elt)
	case *ast.MapType:
		// Map key and value types: map[(<key>)](<value>)
		w.checkParenExpr(n.Key)
		w.checkParenExpr(n.Value)
	case *ast.ChanType:
		// Channel value type: chan (<type>)
		w.checkParenExpr(n.Value)
	case *ast.FuncType:
		// Function parameter and return types
		w.checkFieldList(n.Params)
		w.checkFieldList(n.Results)
	case *ast.TypeAssertExpr:
		// Type assertion: x.(<type>)
		if n.Type != nil {
			w.checkParenExpr(n.Type)
		}
	}
	return w
}

func (w *lintTypeParens) checkFieldList(fl *ast.FieldList) {
	if fl == nil {
		return
	}
	for _, field := range fl.List {
		w.checkParenExpr(field.Type)
	}
}

func (w *lintTypeParens) checkParenExpr(expr ast.Expr) {
	paren, ok := expr.(*ast.ParenExpr)
	if !ok {
		return
	}
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       paren,
		Failure:    "unnecessary parentheses in type expression",
	})
}
