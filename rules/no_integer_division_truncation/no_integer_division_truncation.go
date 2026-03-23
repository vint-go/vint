package no_integer_division_truncation

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/big"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoIntegerDivisionTruncationRule detects integer division of literals where
// the numerator is smaller than the denominator, which always results in zero.
type NoIntegerDivisionTruncationRule struct{}

// Apply applies the rule to given file.
func (r *NoIntegerDivisionTruncationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintIntDivTrunc{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIntegerDivisionTruncationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintIntDivTrunc{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoIntegerDivisionTruncationRule) Name() string {
	return "noIntegerDivisionTruncation"
}

// Group returns the rule group.
func (*NoIntegerDivisionTruncationRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoIntegerDivisionTruncationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIntDivTrunc struct {
	onFailure func(lint.Failure)
}

func (w *lintIntDivTrunc) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.QUO {
		return w
	}

	// Both operands must be integer literals.
	xLit, ok := binExpr.X.(*ast.BasicLit)
	if !ok || xLit.Kind != token.INT {
		return w
	}

	yLit, ok := binExpr.Y.(*ast.BasicLit)
	if !ok || yLit.Kind != token.INT {
		return w
	}

	// Parse the literal values.
	numerator, okN := new(big.Int).SetString(xLit.Value, 0)
	denominator, okD := new(big.Int).SetString(yLit.Value, 0)
	if !okN || !okD {
		return w
	}

	// Skip division by zero (that's a different error).
	if denominator.Sign() == 0 {
		return w
	}

	// Skip zero numerator (0/N is always 0, no truncation).
	if numerator.Sign() == 0 {
		return w
	}

	// Use absolute values for comparison.
	absNum := new(big.Int).Abs(numerator)
	absDen := new(big.Int).Abs(denominator)

	// If |numerator| < |denominator|, the integer division truncates to zero.
	if absNum.Cmp(absDen) < 0 {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       binExpr,
			Failure:    fmt.Sprintf("integer division of %s by %s results in zero; use floating-point division if intended", xLit.Value, yLit.Value),
		})
	}

	return w
}
