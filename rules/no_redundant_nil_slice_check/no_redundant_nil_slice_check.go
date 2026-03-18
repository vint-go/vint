package no_redundant_nil_slice_check

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantNilSliceCheckRule detects redundant nil checks on slices before len() calls.
// For example, `s != nil && len(s) > 0` can be simplified to `len(s) > 0`.
type NoRedundantNilSliceCheckRule struct{}

// Apply applies the rule to the given file.
func (r *NoRedundantNilSliceCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilSliceCheck{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantNilSliceCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilSliceCheck{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantNilSliceCheckRule) Name() string {
	return "noRedundantNilSliceCheck"
}

// Group returns the rule group.
func (*NoRedundantNilSliceCheckRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantNilSliceCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantNilSliceCheck struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantNilSliceCheck) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.LAND {
		return w
	}

	// Check pattern: x != nil && len(x) > 0
	if w.matchNilCheckAndLen(binExpr.X, binExpr.Y) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryStyle,
			Failure:    "redundant nil check on slice; len() handles nil slices",
		})
	}

	return w
}

// matchNilCheckAndLen checks if the two expressions match the pattern:
// nilCheck && lenCheck, where nilCheck is `x != nil` and lenCheck is `len(x) > 0`
// (and variants like `len(x) != 0`, `len(x) >= 1`).
func (w *lintRedundantNilSliceCheck) matchNilCheckAndLen(left, right ast.Expr) bool {
	// left must be: x != nil
	nilIdent := extractNilCheckIdent(left)
	if nilIdent == "" {
		return false
	}

	// right must be: len(x) > 0 or len(x) != 0 or len(x) >= 1
	lenIdent := extractLenCheckIdent(right)
	if lenIdent == "" {
		return false
	}

	return nilIdent == lenIdent
}

// extractNilCheckIdent extracts the variable name from an expression of the form `x != nil`.
func extractNilCheckIdent(expr ast.Expr) string {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return ""
	}

	if binExpr.Op != token.NEQ {
		return ""
	}

	// Check x != nil
	ident, ok := binExpr.X.(*ast.Ident)
	if !ok {
		return ""
	}

	nilIdent, ok := binExpr.Y.(*ast.Ident)
	if !ok {
		return ""
	}

	if nilIdent.Name != "nil" {
		return ""
	}

	return ident.Name
}

// extractLenCheckIdent extracts the variable name from an expression of the form
// `len(x) > 0`, `len(x) != 0`, or `len(x) >= 1`.
func extractLenCheckIdent(expr ast.Expr) string {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return ""
	}

	// Left side must be len(x)
	callExpr, ok := binExpr.X.(*ast.CallExpr)
	if !ok {
		return ""
	}

	funIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return ""
	}

	if funIdent.Name != "len" {
		return ""
	}

	if len(callExpr.Args) != 1 {
		return ""
	}

	argIdent, ok := callExpr.Args[0].(*ast.Ident)
	if !ok {
		return ""
	}

	// Right side must be a literal 0 or 1 depending on operator
	lit, ok := binExpr.Y.(*ast.BasicLit)
	if !ok {
		return ""
	}

	if lit.Kind != token.INT {
		return ""
	}

	switch binExpr.Op {
	case token.GTR: // len(x) > 0
		if lit.Value == "0" {
			return argIdent.Name
		}
	case token.NEQ: // len(x) != 0
		if lit.Value == "0" {
			return argIdent.Name
		}
	case token.GEQ: // len(x) >= 1
		if lit.Value == "1" {
			return argIdent.Name
		}
	}

	return ""
}
