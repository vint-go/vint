package no_redundant_string_concat

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantStringConcatRule detects string concatenation with empty string literals.
type NoRedundantStringConcatRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantStringConcatRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantStringConcat{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantStringConcatRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintRedundantStringConcat{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantStringConcatRule) Name() string {
	return "noRedundantStringConcat"
}

// Group returns the rule group.
func (*NoRedundantStringConcatRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantStringConcatRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantStringConcat struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantStringConcat) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.ADD {
		return w
	}

	if isEmptyStringLit(binExpr.X) || isEmptyStringLit(binExpr.Y) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       binExpr,
			Failure:    "redundant concatenation with empty string",
		})
	}

	return w
}

// isEmptyStringLit checks if the expression is an empty string literal "".
func isEmptyStringLit(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.STRING && lit.Value == `""`
}
