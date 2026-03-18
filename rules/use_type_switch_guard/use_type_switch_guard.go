package use_type_switch_guard

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseTypeSwitchGuardRule detects type switches that can benefit from a type guard
// clause with a variable. It identifies type switch statements that repeatedly
// perform type assertions on the same value within case clauses.
type UseTypeSwitchGuardRule struct{}

// Apply applies the rule to given file.
func (r *UseTypeSwitchGuardRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeSwitchGuard{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTypeSwitchGuardRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeSwitchGuard{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTypeSwitchGuardRule) Name() string {
	return "useTypeSwitchGuard"
}

// Group returns the rule group.
func (*UseTypeSwitchGuardRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeSwitchGuardRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTypeSwitchGuard struct {
	onFailure func(lint.Failure)
}

func (w *lintTypeSwitchGuard) Visit(node ast.Node) ast.Visitor {
	tsStmt, ok := node.(*ast.TypeSwitchStmt)
	if !ok {
		return w
	}

	// Get the variable being switched on and check if there is already a guard assignment.
	switchVar := extractTypeSwitchVar(tsStmt)
	if switchVar == "" {
		// Could not determine the variable being switched; skip.
		return w
	}

	// If the type switch already has a guard variable (switch v := x.(type)),
	// there is nothing to suggest.
	if hasGuardAssignment(tsStmt) {
		return w
	}

	// Check if any single-type case clause performs a type assertion on the switch variable.
	if hasRedundantTypeAssertion(tsStmt, switchVar) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    "use type switch guard (switch v := " + switchVar + ".(type)) to avoid redundant type assertions",
			Node:       tsStmt,
		})
	}

	return w
}

// extractTypeSwitchVar returns the string representation of the variable being
// type-switched. For "switch x.(type)" it returns "x". For "switch x := y.(type)"
// it returns "y".
func extractTypeSwitchVar(ts *ast.TypeSwitchStmt) string {
	if ts.Assign == nil {
		return ""
	}

	switch stmt := ts.Assign.(type) {
	case *ast.ExprStmt:
		// switch x.(type)
		typeAssert, ok := stmt.X.(*ast.TypeAssertExpr)
		if !ok {
			return ""
		}
		return astutils.GoFmt(typeAssert.X)
	case *ast.AssignStmt:
		// switch v := x.(type)
		if len(stmt.Rhs) != 1 {
			return ""
		}
		typeAssert, ok := stmt.Rhs[0].(*ast.TypeAssertExpr)
		if !ok {
			return ""
		}
		return astutils.GoFmt(typeAssert.X)
	}

	return ""
}

// hasGuardAssignment returns true if the type switch has a guard assignment
// (i.e., switch v := x.(type)).
func hasGuardAssignment(ts *ast.TypeSwitchStmt) bool {
	_, ok := ts.Assign.(*ast.AssignStmt)
	return ok
}

// hasRedundantTypeAssertion checks whether any single-type case clause in the
// type switch body contains a type assertion on the given variable.
func hasRedundantTypeAssertion(ts *ast.TypeSwitchStmt, switchVar string) bool {
	if ts.Body == nil {
		return false
	}

	for _, stmt := range ts.Body.List {
		caseClause, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		// Skip case clauses with multiple types (they result in interface{})
		if len(caseClause.List) != 1 {
			continue
		}

		// Check the body for type assertions on switchVar
		if containsTypeAssertionOnVar(caseClause.Body, switchVar) {
			return true
		}
	}

	return false
}

// containsTypeAssertionOnVar returns true if the given list of statements
// contains a type assertion expression on the specified variable.
func containsTypeAssertionOnVar(stmts []ast.Stmt, varName string) bool {
	found := false
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			if found {
				return false
			}
			typeAssert, ok := n.(*ast.TypeAssertExpr)
			if !ok {
				return true
			}
			// Must be a concrete type assertion (not x.(type))
			if typeAssert.Type == nil {
				return true
			}
			if astutils.GoFmt(typeAssert.X) == varName {
				found = true
				return false
			}
			return true
		})
		if found {
			return true
		}
	}
	return found
}
