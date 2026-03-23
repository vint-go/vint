package no_mismatched_append_assign

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMismatchedAppendAssignRule detects suspicious append result assignments where the
// result of append is assigned to a different slice than the one being appended to.
type NoMismatchedAppendAssignRule struct{}

// Apply applies the rule to given file.
func (r *NoMismatchedAppendAssignRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintMismatchedAppendAssign{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMismatchedAppendAssignRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintMismatchedAppendAssign{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMismatchedAppendAssignRule) Name() string {
	return "noMismatchedAppendAssign"
}

// Group returns the rule group.
func (*NoMismatchedAppendAssignRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoMismatchedAppendAssignRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMismatchedAppendAssign struct {
	onFailure func(lint.Failure)
}

func (w *lintMismatchedAppendAssign) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// We only check single-value assignments: x = append(y, ...)
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return w
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if it's a call to the built-in append
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "append" {
		return w
	}

	// Need at least the slice argument
	if len(call.Args) < 1 {
		return w
	}

	lhs := assign.Lhs[0]
	firstArg := call.Args[0]

	// Skip blank identifier assignments: _ = append(x, ...)
	if id, ok := lhs.(*ast.Ident); ok && id.Name == "_" {
		return w
	}

	// Skip if the call uses ellipsis (e.g., xs = append(ys, xs...))
	if call.Ellipsis.IsValid() {
		return w
	}

	// Compare the string representations of LHS and first argument
	lhsStr := astutils.GoFmt(lhs)
	firstArgStr := astutils.GoFmt(firstArg)

	if lhsStr != "" && firstArgStr != "" && lhsStr != firstArgStr {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       assign,
			Category:   lint.FailureCategoryLogic,
			Failure:    "append result is assigned to a different slice '" + lhsStr + "' instead of '" + firstArgStr + "'",
		})
	}

	return w
}
