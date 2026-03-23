package no_nil_func_comparison

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNilFuncComparisonRule detects comparisons of named functions against nil.
// A named function is never nil, so comparing it to nil is always true or always false,
// indicating a likely bug where the programmer intended to call the function.
type NoNilFuncComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoNilFuncComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoNilFuncComparison{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNilFuncComparisonRule) Name() string {
	return "noNilFuncComparison"
}

// Group returns the rule group.
func (*NoNilFuncComparisonRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNilFuncComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoNilFuncComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintNoNilFuncComparison struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoNilFuncComparison) Visit(node ast.Node) ast.Visitor {
	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check == and != comparisons
	if expr.Op != token.EQL && expr.Op != token.NEQ {
		return w
	}

	// Check if one side is nil and the other is a named function
	if isNilIdent(expr.Y) {
		if name := w.namedFuncName(expr.X); name != "" {
			w.reportFailure(expr, name)
		}
	} else if isNilIdent(expr.X) {
		if name := w.namedFuncName(expr.Y); name != "" {
			w.reportFailure(expr, name)
		}
	}

	return w
}

// isNilIdent checks if an expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// namedFuncName checks if an expression refers to a named function (not a function variable)
// and returns its name. Returns empty string if not a named function.
func (w *lintNoNilFuncComparison) namedFuncName(expr ast.Expr) string {
	info := w.pkg.TypesInfo()
	if info == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.Ident:
		obj := info.ObjectOf(e)
		if obj == nil {
			return ""
		}
		if _, ok := obj.(*types.Func); ok {
			return e.Name
		}
	case *ast.SelectorExpr:
		obj := info.ObjectOf(e.Sel)
		if obj == nil {
			return ""
		}
		if _, ok := obj.(*types.Func); ok {
			return fmt.Sprintf("%s.%s", types.ExprString(e.X), e.Sel.Name)
		}
	}

	return ""
}

func (w *lintNoNilFuncComparison) reportFailure(node ast.Node, funcName string) {
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       node,
		Failure:    fmt.Sprintf("comparison of function %s with nil is always the same", funcName),
	})
}
