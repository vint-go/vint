package no_unsigned_negative_comparison

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnsignedNegativeComparisonRule detects pointless comparisons of unsigned
// values against negative values or zero with ordering operators. An unsigned
// integer can never be less than zero, so such comparisons are always true or
// always false and usually indicate a logic error.
type NoUnsignedNegativeComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoUnsignedNegativeComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	w := &lintUnsignedNegCmp{
		typesInfo: typesInfo,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnsignedNegativeComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	w := &lintUnsignedNegCmp{
		typesInfo: typesInfo,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoUnsignedNegativeComparisonRule) Name() string {
	return "noUnsignedNegativeComparison"
}

// Group returns the rule group.
func (*NoUnsignedNegativeComparisonRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnsignedNegativeComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoUnsignedNegativeComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintUnsignedNegCmp struct {
	typesInfo *types.Info
	onFailure func(lint.Failure)
}

func (w *lintUnsignedNegCmp) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check ordering comparisons: <, <=, >, >=
	if !isOrderingOp(binExpr.Op) {
		return w
	}

	// Check both orientations:
	// Case 1: unsigned on left, constant on right
	if result := w.checkPair(binExpr.X, binExpr.Y, binExpr.Op); result != "" {
		w.reportFailure(binExpr, result)
		return w
	}

	// Case 2: constant on left, unsigned on right -- swap the operator
	if result := w.checkPair(binExpr.Y, binExpr.X, swapOp(binExpr.Op)); result != "" {
		w.reportFailure(binExpr, result)
		return w
	}

	return w
}

// checkPair checks if unsignedExpr is unsigned and constExpr is a non-positive
// constant such that the comparison with op is trivially determined.
// Returns "always true", "always false", or "" if not applicable.
func (w *lintUnsignedNegCmp) checkPair(unsignedExpr, constExpr ast.Expr, op token.Token) string {
	if !w.isUnsigned(unsignedExpr) {
		return ""
	}

	constVal := w.constantValue(constExpr)
	if constVal == nil {
		return ""
	}

	sign := constant.Sign(constVal)

	if sign < 0 {
		// Constant is negative: any unsigned value is always > negative.
		switch op {
		case token.LSS, token.LEQ:
			// unsigned < negative  or  unsigned <= negative => always false
			return "always false"
		case token.GTR, token.GEQ:
			// unsigned > negative  or  unsigned >= negative => always true
			return "always true"
		}
	}

	if sign == 0 {
		// Constant is zero.
		switch op {
		case token.LSS:
			// unsigned < 0 => always false
			return "always false"
		case token.GEQ:
			// unsigned >= 0 => always true
			return "always true"
		}
	}

	return ""
}

func (w *lintUnsignedNegCmp) reportFailure(binExpr *ast.BinaryExpr, result string) {
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       binExpr,
		Failure:    fmt.Sprintf("comparing unsigned value against negative value is %s", result),
	})
}

// isUnsigned checks whether an expression has an unsigned integer type.
func (w *lintUnsignedNegCmp) isUnsigned(expr ast.Expr) bool {
	t := w.typesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}

	return basic.Info()&types.IsUnsigned != 0
}

// constantValue returns the constant value of an expression, or nil if not constant.
func (w *lintUnsignedNegCmp) constantValue(expr ast.Expr) constant.Value {
	tv, ok := w.typesInfo.Types[expr]
	if !ok || tv.Value == nil {
		return nil
	}

	if tv.Value.Kind() != constant.Int && tv.Value.Kind() != constant.Float {
		return nil
	}

	return tv.Value
}

// isOrderingOp returns true for <, <=, >, >= operators.
func isOrderingOp(op token.Token) bool {
	switch op {
	case token.LSS, token.LEQ, token.GTR, token.GEQ:
		return true
	}
	return false
}

// swapOp returns the operator with left and right swapped.
func swapOp(op token.Token) token.Token {
	switch op {
	case token.LSS:
		return token.GTR
	case token.LEQ:
		return token.GEQ
	case token.GTR:
		return token.LSS
	case token.GEQ:
		return token.LEQ
	}
	return op
}
