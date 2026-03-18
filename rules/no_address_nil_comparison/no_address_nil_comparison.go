package no_address_nil_comparison

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoAddressNilComparisonRule detects comparisons of address-of expressions against nil.
// The address of a variable is never nil, so such comparisons are always
// false (for ==) or always true (for !=), indicating a logic error.
type NoAddressNilComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoAddressNilComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintAddressNilComparison{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoAddressNilComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintAddressNilComparison{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoAddressNilComparisonRule) Name() string {
	return "noAddressNilComparison"
}

// Group returns the rule group.
func (*NoAddressNilComparisonRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoAddressNilComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintAddressNilComparison struct {
	onFailure func(lint.Failure)
}

func (w *lintAddressNilComparison) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check equality/inequality operators
	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return w
	}

	// Check if one side is an address-of expression and the other is nil
	if isAddressOf(binExpr.X) && isNilIdent(binExpr.Y) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryLogic,
			Failure:    "address of a variable is never nil",
		})
	} else if isNilIdent(binExpr.X) && isAddressOf(binExpr.Y) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryLogic,
			Failure:    "address of a variable is never nil",
		})
	}

	return w
}

// isAddressOf checks if the expression is a unary & (address-of) expression.
func isAddressOf(expr ast.Expr) bool {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	return unary.Op == token.AND
}

// isNilIdent checks if the expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "nil"
}
