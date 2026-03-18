package use_bytes_equal

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseBytesEqualRule detects calls to bytes.Compare(a, b) == 0 that should
// be replaced with bytes.Equal(a, b).
type UseBytesEqualRule struct{}

// Apply applies the rule to given file.
func (r *UseBytesEqualRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintBytesEqual{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseBytesEqualRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintBytesEqual{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseBytesEqualRule) Name() string {
	return "useBytesEqual"
}

// Group returns the rule group.
func (*UseBytesEqualRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseBytesEqualRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBytesEqual struct {
	onFailure func(lint.Failure)
}

func (w *lintBytesEqual) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// We look for patterns like: bytes.Compare(a, b) == 0 or 0 == bytes.Compare(a, b)
	// Also != 0 variants (which could use !bytes.Equal)
	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return w
	}

	if isBytesCompareAgainstZero(binExpr.X, binExpr.Y) || isBytesCompareAgainstZero(binExpr.Y, binExpr.X) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       node,
			Category:   lint.FailureCategoryStyle,
			Failure:    "replace bytes.Compare with bytes.Equal",
		})
	}

	return w
}

// isBytesCompareAgainstZero checks if left is bytes.Compare(...) and right is 0.
func isBytesCompareAgainstZero(left, right ast.Expr) bool {
	call, ok := left.(*ast.CallExpr)
	if !ok {
		return false
	}

	if !astutils.IsPkgDotName(call.Fun, "bytes", "Compare") {
		return false
	}

	return isIntZero(right)
}

// isIntZero checks if expr is the integer literal 0.
func isIntZero(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.INT && lit.Value == "0"
}
