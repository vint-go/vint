package no_imprecise_constant

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoImpreciseConstantRule detects floating-point constants with more precision
// than float64 can represent. Go constants have arbitrary precision, but the
// extra digits are meaningless once assigned to a float variable.
type NoImpreciseConstantRule struct{}

// maxFloat64SignificantDigits is the maximum number of significant decimal
// digits that a float64 can represent.
const maxFloat64SignificantDigits = 17

// Apply applies the rule to given file.
func (r *NoImpreciseConstantRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintImpreciseConst{onFailure: onFailure}

	for _, decl := range file.AST.Decls {
		ast.Walk(w, decl)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImpreciseConstantRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintImpreciseConst{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoImpreciseConstantRule) Name() string {
	return "noImpreciseConstant"
}

// Group returns the rule group.
func (*NoImpreciseConstantRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpreciseConstantRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintImpreciseConst struct {
	onFailure func(lint.Failure)
}

func (w *lintImpreciseConst) Visit(node ast.Node) ast.Visitor {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return w
	}

	if genDecl.Tok != token.CONST {
		return w
	}

	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, val := range vs.Values {
			w.checkExpr(val, vs)
		}
	}

	return w
}

// checkExpr checks whether an expression contains a float literal with
// excessive precision. It handles basic literals and unary expressions
// (e.g., negative constants).
func (w *lintImpreciseConst) checkExpr(expr ast.Expr, vs *ast.ValueSpec) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.FLOAT {
			w.checkFloatLiteral(e, vs)
		}
	case *ast.UnaryExpr:
		// Handle negative constants: const x = -3.14159...
		if lit, ok := e.X.(*ast.BasicLit); ok && lit.Kind == token.FLOAT {
			w.checkFloatLiteral(lit, vs)
		}
	}
}

// checkFloatLiteral checks a single float literal for excessive precision.
func (w *lintImpreciseConst) checkFloatLiteral(lit *ast.BasicLit, vs *ast.ValueSpec) {
	sig := countSignificantDigits(lit.Value)
	if sig > maxFloat64SignificantDigits {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       vs,
			Failure:    "floating-point constant has excessive precision; use math constants or let Go handle the precision",
		})
	}
}

// countSignificantDigits counts the number of significant decimal digits
// in a floating-point literal string. It handles:
//   - Standard notation: "3.14159", "0.001"
//   - Exponent notation: "1.23e10", "1.23E-5"
//   - Underscores in literals: "3.141_592_653"
//   - Hex float literals: "0x1.fp10" (skipped, returns 0)
func countSignificantDigits(s string) int {
	// Remove underscores (Go allows _ in numeric literals)
	s = strings.ReplaceAll(s, "_", "")

	// Skip hex float literals — they don't have the same precision concern
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return 0
	}

	// Strip exponent suffix if present (e.g., "1.23e10" -> "1.23")
	if idx := strings.IndexAny(s, "eE"); idx >= 0 {
		s = s[:idx]
	}

	// Remove sign if present
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}

	// Remove the decimal point to get the digit string
	s = strings.ReplaceAll(s, ".", "")

	// Strip leading zeros (not significant)
	s = strings.TrimLeft(s, "0")

	// Strip trailing zeros only if there was a decimal point in the original
	// (trailing zeros in a float like "1.500" are not extra precision)
	s = strings.TrimRight(s, "0")

	return len(s)
}
