package no_default_slice_index

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDefaultSliceIndexRule detects slice expressions where the high index equals
// len(s), which can be simplified by omitting the high index.
// For example, s[:len(s)] is equivalent to s[:] or just s.
type NoDefaultSliceIndexRule struct{}

// Apply applies the rule to given file.
func (r *NoDefaultSliceIndexRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDefaultSliceIndex{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDefaultSliceIndexRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDefaultSliceIndex{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDefaultSliceIndexRule) Name() string {
	return "noDefaultSliceIndex"
}

// Group returns the rule group.
func (*NoDefaultSliceIndexRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoDefaultSliceIndexRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDefaultSliceIndex struct {
	onFailure func(lint.Failure)
}

func (w *lintDefaultSliceIndex) Visit(node ast.Node) ast.Visitor {
	sliceExpr, ok := node.(*ast.SliceExpr)
	if !ok {
		return w
	}

	// Skip three-index slices (s[:len(s):cap])
	if sliceExpr.Max != nil {
		return w
	}

	// High must be present
	if sliceExpr.High == nil {
		return w
	}

	// Check if High is a call to len(sliceExpr.X)
	if !isLenOfExpr(sliceExpr.High, sliceExpr.X) {
		return w
	}

	fullExprStr := astutils.GoFmt(sliceExpr)

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       sliceExpr,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("omit default slice index: %s can be simplified", fullExprStr),
	})

	return w
}

// isLenOfExpr checks if highExpr is a call to len(x) where x matches sliceX.
func isLenOfExpr(highExpr ast.Expr, sliceX ast.Expr) bool {
	if highExpr == nil {
		return false
	}
	call, ok := highExpr.(*ast.CallExpr)
	if !ok {
		return false
	}
	// Check that the function is the built-in "len"
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "len" {
		return false
	}
	// Must have exactly one argument
	if len(call.Args) != 1 {
		return false
	}
	// The argument must match the slice expression
	return astutils.GoFmt(call.Args[0]) == astutils.GoFmt(sliceX)
}
