package no_yoda_condition

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoYodaConditionRule detects Yoda style conditions and suggests reordering them.
type NoYodaConditionRule struct{}

// Apply applies the rule to given file.
func (r *NoYodaConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoYodaCond{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoYodaConditionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoYodaCond{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoYodaConditionRule) Name() string {
	return "noYodaCondition"
}

// Group returns the rule group.
func (*NoYodaConditionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoYodaConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoYodaCond struct {
	onFailure func(lint.Failure)
}

func (w *lintNoYodaCond) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check comparison operators
	if !isComparisonOp(binExpr.Op) {
		return w
	}

	// Check if left side is a constant/literal and right side is not
	if isConstExpr(binExpr.X) && !isConstExpr(binExpr.Y) {
		lhs := astutils.GoFmt(binExpr.X)
		rhs := astutils.GoFmt(binExpr.Y)
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       binExpr,
			Failure:    fmt.Sprintf("yoda condition: rewrite %s %s %s as %s %s %s", lhs, binExpr.Op.String(), rhs, rhs, swapOp(binExpr.Op), lhs),
		})
	}

	return w
}

// isComparisonOp returns true for comparison operators.
func isComparisonOp(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		return true
	}
	return false
}

// isConstExpr returns true if the expression is a constant or literal value.
func isConstExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// Numeric, string, char literals
		return true
	case *ast.Ident:
		// nil, true, false are constant identifiers
		return e.Name == "nil" || e.Name == "true" || e.Name == "false"
	case *ast.UnaryExpr:
		// Negative numbers like -1
		return isConstExpr(e.X)
	default:
		return false
	}
}

// swapOp returns the reversed comparison operator.
func swapOp(op token.Token) string {
	switch op {
	case token.LSS:
		return ">"
	case token.LEQ:
		return ">="
	case token.GTR:
		return "<"
	case token.GEQ:
		return "<="
	default:
		return op.String() // EQL and NEQ are symmetric
	}
}
