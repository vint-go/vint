package use_direct_string_comparison

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseDirectStringComparisonRule detects empty string checks using len(s) == 0
// that can be written more idiomatically as s == "".
type UseDirectStringComparisonRule struct{}

// Apply applies the rule to given file.
func (r *UseDirectStringComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDirectStringComparison{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseDirectStringComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDirectStringComparison{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseDirectStringComparisonRule) Name() string {
	return "useDirectStringComparison"
}

// Group returns the rule group.
func (*UseDirectStringComparisonRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseDirectStringComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDirectStringComparison struct {
	onFailure func(lint.Failure)
}

func (w *lintDirectStringComparison) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return w
	}

	// Check for len(s) == 0 or len(s) != 0
	if isLenCall(binExpr.X) && isZeroLiteral(binExpr.Y) {
		w.report(binExpr)
		return w
	}

	// Check for 0 == len(s) or 0 != len(s)
	if isZeroLiteral(binExpr.X) && isLenCall(binExpr.Y) {
		w.report(binExpr)
		return w
	}

	return w
}

func (w *lintDirectStringComparison) report(node ast.Node) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryStyle,
		Failure:    "use direct string comparison instead of len() comparison",
	})
}

// isLenCall checks if the expression is a call to the builtin len function.
func isLenCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "len" && len(call.Args) == 1
}

// isZeroLiteral checks if the expression is the integer literal 0.
func isZeroLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}

	return lit.Kind == token.INT && lit.Value == "0"
}
