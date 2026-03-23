package use_type_assert_result

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTypeAssertResultRule detects type switches where the switched variable is
// re-asserted inside case clauses. In a type switch, you can capture the typed
// value with a guard variable (switch v := x.(type)) instead of re-asserting.
type UseTypeAssertResultRule struct{}

// Apply applies the rule to given file.
func (r *UseTypeAssertResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeAssertResult{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTypeAssertResultRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeAssertResult{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTypeAssertResultRule) Name() string {
	return "useTypeAssertResult"
}

// Group returns the rule group.
func (*UseTypeAssertResultRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeAssertResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTypeAssertResult struct {
	onFailure func(lint.Failure)
}

func (w *lintTypeAssertResult) Visit(node ast.Node) ast.Visitor {
	tsStmt, ok := node.(*ast.TypeSwitchStmt)
	if !ok {
		return w
	}

	// Get the variable being switched on.
	switchVar := extractTypeSwitchVar(tsStmt)
	if switchVar == "" {
		return w
	}

	// If the type switch already has a guard variable (switch v := x.(type)),
	// the variable is already captured; no redundant assertions to flag.
	if hasGuardAssignment(tsStmt) {
		return w
	}

	// Check each single-type case clause for type assertions on switchVar.
	if tsStmt.Body != nil {
		for _, stmt := range tsStmt.Body.List {
			caseClause, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}
			// Skip case clauses with multiple types or default (no types).
			if len(caseClause.List) != 1 {
				continue
			}
			// Find and report each redundant type assertion in the case body.
			reportRedundantAssertions(caseClause.Body, switchVar, w.onFailure)
		}
	}

	return w
}

// extractTypeSwitchVar returns the string representation of the variable being
// type-switched. For "switch x.(type)" it returns "x".
func extractTypeSwitchVar(ts *ast.TypeSwitchStmt) string {
	if ts.Assign == nil {
		return ""
	}

	exprStmt, ok := ts.Assign.(*ast.ExprStmt)
	if !ok {
		return ""
	}
	typeAssert, ok := exprStmt.X.(*ast.TypeAssertExpr)
	if !ok {
		return ""
	}
	return astutils.GoFmt(typeAssert.X)
}

// hasGuardAssignment returns true if the type switch has a guard assignment
// (i.e., switch v := x.(type)).
func hasGuardAssignment(ts *ast.TypeSwitchStmt) bool {
	_, ok := ts.Assign.(*ast.AssignStmt)
	return ok
}

// reportRedundantAssertions walks statements and reports each type assertion
// on the given variable.
func reportRedundantAssertions(stmts []ast.Stmt, varName string, onFailure func(lint.Failure)) {
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			typeAssert, ok := n.(*ast.TypeAssertExpr)
			if !ok {
				return true
			}
			// Must be a concrete type assertion (not x.(type))
			if typeAssert.Type == nil {
				return true
			}
			if astutils.GoFmt(typeAssert.X) == varName {
				typeName := astutils.GoFmt(typeAssert.Type)
				onFailure(lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryStyle,
					Failure:    "use the result of the type switch instead of asserting " + varName + ".(" + typeName + ")",
					Node:       typeAssert,
				})
			}
			return true
		})
	}
}
