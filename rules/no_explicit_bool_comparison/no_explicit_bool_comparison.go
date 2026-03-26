package no_explicit_bool_comparison

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExplicitBoolComparisonRule detects explicit comparisons with boolean constants.
// Comparing a boolean value to true or false is redundant.
type NoExplicitBoolComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoExplicitBoolComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.IsTest() {
		return nil
	}

	var failures []lint.Failure

	w := &lintNoExplicitBoolComparison{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExplicitBoolComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if file.IsTest() {
		return nil
	}

	var failures []lint.Failure

	w := &lintNoExplicitBoolComparison{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoExplicitBoolComparisonRule) Name() string {
	return "noExplicitBoolComparison"
}

// Group returns the rule group.
func (*NoExplicitBoolComparisonRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoExplicitBoolComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoExplicitBoolComparison struct {
	onFailure func(lint.Failure)
}

func (w *lintNoExplicitBoolComparison) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Only check == and != operators
	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return w
	}

	// Check if either side is a boolean literal (true or false)
	boolLit, boolSide := getBoolLiteral(binExpr)
	if boolLit == "" {
		return w
	}

	// Get the other expression (non-boolean side)
	var otherExpr ast.Expr
	if boolSide == "left" {
		otherExpr = binExpr.Y
	} else {
		otherExpr = binExpr.X
	}

	otherStr := astutils.GoFmt(otherExpr)

	// Determine the suggested replacement
	var suggestion string
	if (binExpr.Op == token.EQL && boolLit == "true") || (binExpr.Op == token.NEQ && boolLit == "false") {
		suggestion = otherStr
	} else {
		suggestion = "!" + otherStr
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       binExpr,
		Failure:    fmt.Sprintf("omit explicit comparison to boolean constant, can be simplified to %s", suggestion),
	})

	return w
}

// getBoolLiteral checks if either side of a binary expression is a boolean literal.
// Returns the literal value ("true" or "false") and which side ("left" or "right").
func getBoolLiteral(binExpr *ast.BinaryExpr) (string, string) {
	if ident, ok := binExpr.X.(*ast.Ident); ok {
		if ident.Name == "true" || ident.Name == "false" {
			return ident.Name, "left"
		}
	}
	if ident, ok := binExpr.Y.(*ast.Ident); ok {
		if ident.Name == "true" || ident.Name == "false" {
			return ident.Name, "right"
		}
	}
	return "", ""
}
