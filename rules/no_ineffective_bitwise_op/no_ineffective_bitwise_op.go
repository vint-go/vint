package no_ineffective_bitwise_op

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIneffectiveBitwiseOpRule detects bitwise operations that have no effect.
type NoIneffectiveBitwiseOpRule struct{}

// Apply applies the rule to given file.
func (r *NoIneffectiveBitwiseOpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoIneffectiveBitwiseOp{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIneffectiveBitwiseOpRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoIneffectiveBitwiseOp{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoIneffectiveBitwiseOpRule) Name() string {
	return "noIneffectiveBitwiseOp"
}

// Group returns the rule group.
func (*NoIneffectiveBitwiseOpRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoIneffectiveBitwiseOpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoIneffectiveBitwiseOp struct {
	onFailure func(lint.Failure)
}

func (w *lintNoIneffectiveBitwiseOp) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if !isBitwiseOp(binExpr.Op) {
		return w
	}

	// Check if either operand is the integer literal 0
	lhsIsZero := isIntLitZero(binExpr.X)
	rhsIsZero := isIntLitZero(binExpr.Y)

	if !lhsIsZero && !rhsIsZero {
		return w
	}

	var msg string
	switch binExpr.Op {
	case token.XOR:
		// x ^ 0 == x, 0 ^ x == x: ineffective, always equals the other operand
		msg = "ineffective bitwise XOR with 0: the expression always equals the other operand"
	case token.OR:
		// x | 0 == x, 0 | x == x: ineffective, always equals the other operand
		msg = "ineffective bitwise OR with 0: the expression always equals the other operand"
	case token.AND:
		// x & 0 == 0, 0 & x == 0: ineffective, always equals 0
		msg = "ineffective bitwise AND with 0: the expression always equals 0"
	case token.AND_NOT:
		if rhsIsZero {
			// x &^ 0 == x: ineffective, always equals x
			msg = "ineffective bitwise AND NOT with 0: the expression always equals the other operand"
		} else if lhsIsZero {
			// 0 &^ x == 0: ineffective, always equals 0
			msg = "ineffective bitwise AND NOT with 0: the expression always equals 0"
		}
	}

	if msg != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       binExpr,
			Failure:    msg,
		})
	}

	return w
}

// isBitwiseOp returns true for bitwise operators.
func isBitwiseOp(op token.Token) bool {
	switch op {
	case token.AND, token.OR, token.XOR, token.AND_NOT:
		return true
	}
	return false
}

// isIntLitZero checks whether an expression is the integer literal 0.
func isIntLitZero(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	if lit.Kind != token.INT {
		return false
	}
	return lit.Value == "0"
}
