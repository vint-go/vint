package no_range_append_all

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRangeAppendAllRule detects appending an entire slice inside a range loop
// over that same slice, which leads to quadratic behavior.
type NoRangeAppendAllRule struct{}

// Apply applies the rule to given file.
func (r *NoRangeAppendAllRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRangeAppendAll{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRangeAppendAllRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRangeAppendAll{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRangeAppendAllRule) Name() string {
	return "noRangeAppendAll"
}

// Group returns the rule group.
func (*NoRangeAppendAllRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoRangeAppendAllRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRangeAppendAll struct {
	onFailure func(lint.Failure)
}

func (w *lintRangeAppendAll) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Get the string representation of the slice being ranged over
	rangeXStr := astutils.GoFmt(rangeStmt.X)
	if rangeXStr == "" {
		return w
	}

	// Walk the body looking for append calls with ellipsis that use the same slice
	bodyWalker := &rangeAppendAllBodyWalker{
		rangeXStr: rangeXStr,
		onFailure: w.onFailure,
	}
	ast.Walk(bodyWalker, rangeStmt.Body)

	return w
}

type rangeAppendAllBodyWalker struct {
	rangeXStr string
	onFailure func(lint.Failure)
}

func (bw *rangeAppendAllBodyWalker) Visit(node ast.Node) ast.Visitor {
	// Check for assignment statements: x = append(x, y...)
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return bw
	}

	if len(assign.Rhs) != 1 {
		return bw
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return bw
	}

	// Check if it's a call to the built-in append
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "append" {
		return bw
	}

	// Must have the ellipsis operator (variadic expansion)
	if !call.Ellipsis.IsValid() {
		return bw
	}

	// Need at least 2 args for append(slice, otherSlice...)
	if len(call.Args) < 2 {
		return bw
	}

	// Check if the ellipsis argument is the same slice being ranged over
	lastArg := call.Args[len(call.Args)-1]
	lastArgStr := astutils.GoFmt(lastArg)

	if lastArgStr == bw.rangeXStr {
		bw.onFailure(lint.Failure{
			Confidence: 1,
			Node:       assign,
			Category:   lint.FailureCategoryLogic,
			Failure:    fmt.Sprintf("appending entire slice '%s' inside range loop over '%s' causes quadratic behavior", lastArgStr, bw.rangeXStr),
		})
	}

	return bw
}
