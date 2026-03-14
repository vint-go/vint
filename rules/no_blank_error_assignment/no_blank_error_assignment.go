package no_blank_error_assignment

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/lint"
)

// NoBlankErrorAssignmentRule warns when error return values are explicitly assigned to the blank identifier.
type NoBlankErrorAssignmentRule struct{}

// Apply applies the rule to given file.
func (r *NoBlankErrorAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoBlankErrorAssignment{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoBlankErrorAssignmentRule) Name() string {
	return "noBlankErrorAssignment"
}

// Group returns the rule group.
func (*NoBlankErrorAssignmentRule) Group() string {
	return "correctness"
}

type lintNoBlankErrorAssignment struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoBlankErrorAssignment) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Find which LHS positions are blank identifiers
	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name != "_" {
			continue
		}

		// Check if the type assigned to this blank identifier is error
		if w.isErrorAtIndex(assign, i) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryErrors,
				Confidence: 1,
				Node:       assign,
				Failure:    "error assigned to blank identifier",
			})
			// Only report once per assignment statement
			return w
		}
	}

	return w
}

func (w *lintNoBlankErrorAssignment) isErrorAtIndex(assign *ast.AssignStmt, index int) bool {
	// For single-value RHS with multiple LHS (tuple assignment like `a, b := f()`),
	// check the function return types
	if len(assign.Rhs) == 1 && len(assign.Lhs) > 1 {
		return w.isTupleErrorAtIndex(assign.Rhs[0], index)
	}

	// For matching LHS/RHS count (e.g., `_ = f.Close()`),
	// check the type of the corresponding RHS expression
	if index < len(assign.Rhs) {
		t := w.pkg.TypeOf(assign.Rhs[index])
		return t != nil && isErrorType(t)
	}

	return false
}

func (w *lintNoBlankErrorAssignment) isTupleErrorAtIndex(expr ast.Expr, index int) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	tuple, ok := t.(*types.Tuple)
	if !ok {
		return false
	}

	if index >= tuple.Len() {
		return false
	}

	return isErrorType(tuple.At(index).Type())
}

func isErrorType(t types.Type) bool {
	// Check if the type implements the error interface
	// The error interface is a named type with Id "_.error"
	named, ok := t.(*types.Named)
	if ok {
		return named.Obj().Id() == "_.error"
	}

	return false
}
