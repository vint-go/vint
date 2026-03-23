package use_slice_append

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseSliceAppendRule detects for-range loops that append each element from one
// slice to another and can be replaced with a single append using the ... operator.
type UseSliceAppendRule struct{}

// Apply applies the rule to given file.
func (r *UseSliceAppendRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSliceAppend{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSliceAppendRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSliceAppend{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseSliceAppendRule) Name() string {
	return "useSliceAppend"
}

// Group returns the rule group.
func (*UseSliceAppendRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseSliceAppendRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSliceAppend struct {
	onFailure func(lint.Failure)
}

func (w *lintSliceAppend) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// We need: for _, v := range src { dst = append(dst, v) }
	// The body must contain exactly one statement.
	if rangeStmt.Body == nil || len(rangeStmt.Body.List) != 1 {
		return w
	}

	// The range value variable must be present and not blank.
	if rangeStmt.Value == nil {
		return w
	}
	valIdent, ok := rangeStmt.Value.(*ast.Ident)
	if !ok || valIdent.Name == "_" {
		return w
	}

	// The single statement must be an assignment: dst = append(dst, v)
	assignStmt, ok := rangeStmt.Body.List[0].(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Must be a simple assignment with one LHS and one RHS.
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return w
	}

	// RHS must be a call to append.
	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	funIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok || funIdent.Name != "append" {
		return w
	}

	// append must have exactly 2 arguments and no ellipsis.
	if len(callExpr.Args) != 2 || callExpr.Ellipsis.IsValid() {
		return w
	}

	// The first argument to append must match the LHS of the assignment.
	lhsStr := astutils.GoFmt(assignStmt.Lhs[0])
	firstArgStr := astutils.GoFmt(callExpr.Args[0])
	if lhsStr == "" || lhsStr != firstArgStr {
		return w
	}

	// The second argument to append must be the range value variable.
	secondArgStr := astutils.GoFmt(callExpr.Args[1])
	if secondArgStr != valIdent.Name {
		return w
	}

	// Get the source slice being ranged over.
	srcStr := astutils.GoFmt(rangeStmt.X)

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       rangeStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("use '%s = append(%s, %s...)' instead of a loop to append slice elements", lhsStr, lhsStr, srcStr),
	})

	return w
}
