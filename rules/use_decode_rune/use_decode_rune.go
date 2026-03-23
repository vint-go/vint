package use_decode_rune

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseDecodeRuneRule detects expressions like []rune(s)[0] that perform an
// unnecessary full string-to-rune-slice conversion when only the first rune
// is needed. Using utf8.DecodeRuneInString(s) is more efficient.
type UseDecodeRuneRule struct{}

// Apply applies the rule to given file.
func (r *UseDecodeRuneRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDecodeRune{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseDecodeRuneRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDecodeRune{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseDecodeRuneRule) Name() string {
	return "useDecodeRune"
}

// Group returns the rule group.
func (*UseDecodeRuneRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseDecodeRuneRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDecodeRune struct {
	onFailure func(lint.Failure)
}

func (w *lintDecodeRune) Visit(node ast.Node) ast.Visitor {
	indexExpr, ok := node.(*ast.IndexExpr)
	if !ok {
		return w
	}

	// Check that the index is the integer literal 0
	indexLit, ok := indexExpr.Index.(*ast.BasicLit)
	if !ok || indexLit.Kind != token.INT || indexLit.Value != "0" {
		return w
	}

	// Check that the indexed expression is a call expression (type conversion)
	callExpr, ok := indexExpr.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check that the call is a type conversion to []rune
	if !isRuneSliceType(callExpr.Fun) {
		return w
	}

	// Must have exactly one argument (the string being converted)
	if len(callExpr.Args) != 1 {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryOptimization,
		Failure:    "use utf8.DecodeRuneInString instead of []rune(s)[0] to avoid full string-to-rune-slice conversion",
	})

	return w
}

// isRuneSliceType checks if the expression represents the type []rune.
func isRuneSliceType(expr ast.Expr) bool {
	arrayType, ok := expr.(*ast.ArrayType)
	if !ok {
		return false
	}

	// Must be a slice (no length specified)
	if arrayType.Len != nil {
		return false
	}

	// The element type must be "rune"
	ident, ok := arrayType.Elt.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "rune"
}
