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
	case *ast.ValueSpec:
		// var/const declarations: var x *(int)
		w.checkTypeExpr(n.Type)
	case *ast.TypeSpec:
		// type declarations: type T *(int)
		w.checkTypeExpr(n.Type)
	case *ast.CompositeLit:
		// composite literals: *(int){...}
		w.checkTypeExpr(n.Type)
	case *ast.StructType:
		// struct field types
		w.checkFieldList(n.Fields)
	case *ast.InterfaceType:
		// interface method types
		w.checkFieldList(n.Methods)
	case *ast.ArrayType:
		// Array/slice element type: [N](<type>) or [](<type>)
		w.checkTypeExpr(n.Elt)
	case *ast.MapType:
		// Map key and value types: map[(<key>)](<value>)
		w.checkTypeExpr(n.Key)
		w.checkTypeExpr(n.Value)
	case *ast.ChanType:
		// Channel value type: chan (<type>)
		w.checkTypeExpr(n.Value)
	case *ast.FuncType:
		// Function parameter and return types
		w.checkFieldList(n.Params)
		w.checkFieldList(n.Results)
	case *ast.TypeAssertExpr:
		// Type assertion: x.(<type>)
		if n.Type != nil {
			w.checkTypeExpr(n.Type)
		}
	}
	return w
}

func (w *lintTypeParens) checkFieldList(fl *ast.FieldList) {
	if fl == nil {
		return
	}
	for _, field := range fl.List {
		w.checkTypeExpr(field.Type)
	}
}

// checkTypeExpr checks a type expression for unnecessary parentheses.
// It handles both direct ParenExpr (e.g., (int) in map[(string)]int)
// and StarExpr wrapping ParenExpr (e.g., *(int) in var x *(int)).
func (w *lintTypeParens) checkTypeExpr(expr ast.Expr) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *ast.ParenExpr:
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       e,
			Failure:    "unnecessary parentheses in type expression",
		})
	case *ast.StarExpr:
		// In a type context, *ast.StarExpr is a pointer type.
		// Check if the inner type has unnecessary parens: *(int) -> *int
		if paren, ok := e.X.(*ast.ParenExpr); ok {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       paren,
				Failure:    "unnecessary parentheses in type expression",
			})
		}
	}
}
