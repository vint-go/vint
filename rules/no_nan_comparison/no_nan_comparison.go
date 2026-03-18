package no_nan_comparison

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNanComparisonRule detects comparisons with math.NaN(), which are always
// false because NaN is not equal to anything, including itself.
type NoNanComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoNanComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNanComparison{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNanComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNanComparison{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNanComparisonRule) Name() string {
	return "noNanComparison"
}

// Group returns the rule group.
func (*NoNanComparisonRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNanComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNanComparison struct {
	onFailure func(lint.Failure)
}

func (w *lintNanComparison) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check comparison operators
	if !isComparisonOp(binExpr.Op) {
		return w
	}

	// Check if either side is a call to math.NaN()
	if isNaNCall(binExpr.X) || isNaNCall(binExpr.Y) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryLogic,
			Failure:    "comparing with math.NaN() is always false, use math.IsNaN() instead",
		})
	}

	return w
}

func isComparisonOp(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return true
	}
	return false
}

func isNaNCall(expr ast.Expr) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	return astutils.IsPkgDotName(callExpr.Fun, "math", "NaN")
}
