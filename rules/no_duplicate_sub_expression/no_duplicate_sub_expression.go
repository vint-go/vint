package no_duplicate_sub_expression

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDuplicateSubExpressionRule detects suspicious duplicated sub-expressions in binary expressions.
type NoDuplicateSubExpressionRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateSubExpressionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDupSubExpr{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateSubExpressionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDupSubExpr{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateSubExpressionRule) Name() string {
	return "noDuplicateSubExpression"
}

// Group returns the rule group.
func (*NoDuplicateSubExpressionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateSubExpressionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoDupSubExpr struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDupSubExpr) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Skip operators where duplicate operands are intentional or harmless.
	// For example, x + x is equivalent to 2*x and may be intentional.
	// We focus on comparison, logical, and bitwise operators where duplicates
	// are almost always bugs.
	if !isSuspiciousOp(binExpr.Op) {
		return w
	}

	lhs := astutils.GoFmt(binExpr.X)
	rhs := astutils.GoFmt(binExpr.Y)

	if lhs == "" || rhs == "" {
		return w
	}

	if lhs == rhs {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       binExpr,
			Failure:    "suspicious identical sub-expressions on both sides of '" + binExpr.Op.String() + "'",
		})
	}

	return w
}

// isSuspiciousOp returns true for operators where having the same sub-expression
// on both sides is suspicious and likely a bug.
func isSuspiciousOp(op token.Token) bool {
	switch op {
	case
		// Logical operators: x && x, x || x
		token.LAND, token.LOR,
		// Comparison operators: x == x, x != x, x < x, x <= x, x > x, x >= x
		token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ,
		// Bitwise operators: x & x, x | x, x ^ x, x &^ x
		token.AND, token.OR, token.XOR, token.AND_NOT,
		// Arithmetic operators where same operands are suspicious:
		// x - x (always 0), x / x (always 1), x % x (always 0)
		token.SUB, token.QUO, token.REM:
		return true
	}
	return false
}
