package no_unchecked_type_assertion

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/lint"
)

// NoUncheckedTypeAssertionRule detects when type assertion results are not checked for success.
type NoUncheckedTypeAssertionRule struct{}

// Apply applies the rule to given file.
func (r *NoUncheckedTypeAssertionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect all type assertion expressions that are in safe 2-value assignments
	safeAssertions := map[*ast.TypeAssertExpr]bool{}

	// First pass: find all safe type assertions (comma-ok idiom with non-blank ok)
	ast.Inspect(file.AST, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		if len(assign.Rhs) == 0 {
			return true
		}
		e, ok := assign.Rhs[0].(*ast.TypeAssertExpr)
		if !ok || e == nil {
			return true
		}
		// Type switch uses nil Type
		if e.Type == nil {
			safeAssertions[e] = true
			return true
		}
		// Two LHS where ok value is captured (not _) is safe
		if len(assign.Lhs) == 2 && !astutils.IsIdent(assign.Lhs[1], "_") {
			safeAssertions[e] = true
		}
		return true
	})

	// Second pass: find all TypeAssertExpr that are not safe
	ast.Inspect(file.AST, func(n ast.Node) bool {
		e, ok := n.(*ast.TypeAssertExpr)
		if !ok {
			return true
		}
		// Type switch (Type is nil)
		if e.Type == nil {
			return true
		}
		// Skip if it's in a safe assignment
		if safeAssertions[e] {
			return true
		}

		s := fmt.Sprintf("type assertion result is unchecked in %v, use the comma-ok idiom", astutils.GoFmt(e))
		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       e,
			Failure:    s,
		})

		return true
	})

	return failures
}

// Name returns the rule name.
func (*NoUncheckedTypeAssertionRule) Name() string {
	return "noUncheckedTypeAssertion"
}

// Group returns the rule group.
func (*NoUncheckedTypeAssertionRule) Group() string {
	return "correctness"
}
