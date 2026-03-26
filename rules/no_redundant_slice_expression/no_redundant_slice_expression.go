package no_redundant_slice_expression

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantSliceExpressionRule detects slice expressions that can be simplified
// to the expression itself, such as s[:] or s[0:len(s)].
type NoRedundantSliceExpressionRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantSliceExpressionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantSlice{onFailure: onFailure, typesInfo: typesInfo}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantSliceExpressionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantSlice{onFailure: onFailure, typesInfo: typesInfo}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantSliceExpressionRule) Name() string {
	return "noRedundantSliceExpression"
}

// Group returns the rule group.
func (*NoRedundantSliceExpressionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantSliceExpressionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoRedundantSliceExpressionRule) RequiresTypecheck() bool {
	return true
}

type lintRedundantSlice struct {
	onFailure func(lint.Failure)
	typesInfo *types.Info
}

func (w *lintRedundantSlice) Visit(node ast.Node) ast.Visitor {
	sliceExpr, ok := node.(*ast.SliceExpr)
	if !ok {
		return w
	}

	// Skip array types: array[:] is a necessary conversion from [N]T to []T.
	if w.isArrayType(sliceExpr.X) {
		return w
	}

	// Check for s[:] — both Low and High are nil
	if sliceExpr.Low == nil && sliceExpr.High == nil && sliceExpr.Max == nil {
		xStr := astutils.GoFmt(sliceExpr.X)
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       sliceExpr,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("redundant slice expression %s[:] can be simplified to %s", xStr, xStr),
		})
		return w
	}

	// Check for s[0:len(s)] — Low is 0, High is len(s)
	if sliceExpr.Max == nil && isZeroLiteral(sliceExpr.Low) && isLenOfExpr(sliceExpr.High, sliceExpr.X) {
		xStr := astutils.GoFmt(sliceExpr.X)
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       sliceExpr,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("redundant slice expression %s[0:len(%s)] can be simplified to %s", xStr, xStr, xStr),
		})
		return w
	}

	return w
}

// isZeroLiteral checks if the expression is the integer literal 0.
func isZeroLiteral(expr ast.Expr) bool {
	if expr == nil {
		return false
	}
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.INT && lit.Value == "0"
}

// isArrayType returns true if the expression has an array type. If type info
// is unavailable, it returns true (erring on the safe side by not flagging).
func (w *lintRedundantSlice) isArrayType(expr ast.Expr) bool {
	t := w.typesInfo.TypeOf(expr)
	if t == nil {
		return true // err on safe side
	}
	_, isArray := t.Underlying().(*types.Array)
	return isArray
}

// isLenOfExpr checks if highExpr is a call to len(x) where x matches sliceX.
func isLenOfExpr(highExpr ast.Expr, sliceX ast.Expr) bool {
	if highExpr == nil {
		return false
	}
	call, ok := highExpr.(*ast.CallExpr)
	if !ok {
		return false
	}
	// Check that the function is the built-in "len"
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "len" {
		return false
	}
	// Must have exactly one argument
	if len(call.Args) != 1 {
		return false
	}
	// The argument must match the slice expression
	return astutils.GoFmt(call.Args[0]) == astutils.GoFmt(sliceX)
}
