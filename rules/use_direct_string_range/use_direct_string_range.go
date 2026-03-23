package use_direct_string_range

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseDirectStringRangeRule detects unnecessary conversion to []rune before ranging over a string.
type UseDirectStringRangeRule struct{}

// Apply applies the rule to given file.
func (r *UseDirectStringRangeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDirectStringRange{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseDirectStringRangeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDirectStringRange{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseDirectStringRangeRule) Name() string {
	return "useDirectStringRange"
}

// Group returns the rule group.
func (*UseDirectStringRangeRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseDirectStringRangeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDirectStringRange struct {
	onFailure func(lint.Failure)
}

func (w *lintDirectStringRange) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Check if the range expression is a []rune(...) conversion call
	if isRuneSliceConversion(rangeStmt.X) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       rangeStmt.X,
			Category:   lint.FailureCategoryStyle,
			Failure:    "range over string directly instead of converting to []rune",
		})
	}

	return w
}

// isRuneSliceConversion checks if the expression is a []rune(x) type conversion.
func isRuneSliceConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// The function part should be []rune, which is an *ast.ArrayType
	arrayType, ok := call.Fun.(*ast.ArrayType)
	if !ok {
		return false
	}

	// The element type should be "rune"
	ident, ok := arrayType.Elt.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "rune"
}
