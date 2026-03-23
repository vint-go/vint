package no_pointer_to_ref_param

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoPointerToRefParamRule detects function parameters and return types
// that are pointers to reference types (maps, channels, interfaces).
type NoPointerToRefParamRule struct{}

// Apply applies the rule to the given file.
func (r *NoPointerToRefParamRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintPointerToRef{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoPointerToRefParamRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintPointerToRef{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoPointerToRefParamRule) Name() string {
	return "noPointerToRefParam"
}

// Group returns the rule group.
func (*NoPointerToRefParamRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoPointerToRefParamRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintPointerToRef struct {
	onFailure func(lint.Failure)
}

func (w *lintPointerToRef) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFuncType(n.Type)
	case *ast.FuncLit:
		w.checkFuncType(n.Type)
	}
	return w
}

func (w *lintPointerToRef) checkFuncType(ft *ast.FuncType) {
	if ft.Params != nil {
		w.checkFieldList(ft.Params, "input")
	}
	if ft.Results != nil {
		w.checkFieldList(ft.Results, "output")
	}
}

func (w *lintPointerToRef) checkFieldList(fl *ast.FieldList, direction string) {
	for _, field := range fl.List {
		w.checkType(field.Type, direction)
	}
}

func (w *lintPointerToRef) checkType(expr ast.Expr, direction string) {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return
	}

	refType := refTypeName(star.X)
	if refType == "" {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       expr,
		Failure:    fmt.Sprintf("%s parameter should not be a pointer to %s", direction, refType),
	})
}

// refTypeName returns the name of the reference type if the expression
// represents a map, channel, or interface type. Returns "" otherwise.
func refTypeName(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.MapType:
		return "map"
	case *ast.ChanType:
		return "chan"
	case *ast.InterfaceType:
		return "interface"
	}
	return ""
}
