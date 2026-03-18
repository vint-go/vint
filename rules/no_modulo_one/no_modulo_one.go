package no_modulo_one

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoModuloOneRule detects modulo operations with a divisor of 1, which always yield zero.
type NoModuloOneRule struct{}

// Apply applies the rule to given file.
func (r *NoModuloOneRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoModOne{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoModuloOneRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoModOne{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoModuloOneRule) Name() string {
	return "noModuloOne"
}

// Group returns the rule group.
func (*NoModuloOneRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoModuloOneRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoModOne struct {
	onFailure func(lint.Failure)
}

func (w *lintNoModOne) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.REM {
		return w
	}

	// Check if the right-hand side is the integer literal 1
	lit, ok := binExpr.Y.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind == token.INT && lit.Value == "1" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       binExpr,
			Failure:    "x % 1 is always zero",
		})
	}

	return w
}
