package no_redundant_nil_type_check

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantNilTypeCheckRule detects redundant nil checks before type assertions.
// A two-value type assertion already returns false for nil interfaces, so checking
// x != nil before a type assertion is redundant.
type NoRedundantNilTypeCheckRule struct{}

// Apply applies the rule to the given file.
func (r *NoRedundantNilTypeCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilTypeCheck{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantNilTypeCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilTypeCheck{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantNilTypeCheckRule) Name() string {
	return "noRedundantNilTypeCheck"
}

// Group returns the rule group.
func (*NoRedundantNilTypeCheckRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantNilTypeCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantNilTypeCheck struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantNilTypeCheck) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// The outer if must have no init statement and no else branch
	if ifStmt.Init != nil || ifStmt.Else != nil {
		return w
	}

	// The condition must be: x != nil (or nil != x)
	nilCheckVar := extractNilNeqCheckIdent(ifStmt.Cond)
	if nilCheckVar == "" {
		return w
	}

	// The body must contain exactly one statement
	if len(ifStmt.Body.List) != 1 {
		return w
	}

	// That statement must be an if with a two-value type assertion
	innerIf, ok := ifStmt.Body.List[0].(*ast.IfStmt)
	if !ok {
		return w
	}

	// The inner if must have an init statement that is a type assertion assignment
	assertVar := extractTypeAssertVar(innerIf)
	if assertVar == "" {
		return w
	}

	// The nil-checked variable must be the same as the type-asserted variable
	if nilCheckVar != assertVar {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ifStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    "redundant nil check on variable before type assertion; type assertion already handles nil",
	})

	return w
}

// extractNilNeqCheckIdent extracts the variable name from an expression
// of the form `x != nil` or `nil != x`.
func extractNilNeqCheckIdent(expr ast.Expr) string {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return ""
	}

	if binExpr.Op != token.NEQ {
		return ""
	}

	// Check x != nil
	if ident, ok := binExpr.X.(*ast.Ident); ok {
		if nilIdent, ok := binExpr.Y.(*ast.Ident); ok && nilIdent.Name == "nil" {
			return ident.Name
		}
	}

	// Check nil != x
	if nilIdent, ok := binExpr.X.(*ast.Ident); ok && nilIdent.Name == "nil" {
		if ident, ok := binExpr.Y.(*ast.Ident); ok {
			return ident.Name
		}
	}

	return ""
}

// extractTypeAssertVar extracts the variable name being type-asserted in
// an if statement of the form: if v, ok := x.(T); ok { ... }
func extractTypeAssertVar(ifStmt *ast.IfStmt) string {
	if ifStmt.Init == nil {
		return ""
	}

	// The init must be an assignment: v, ok := x.(T)
	assignStmt, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return ""
	}

	// Must have 2 LHS values (v, ok)
	if len(assignStmt.Lhs) != 2 || len(assignStmt.Rhs) != 1 {
		return ""
	}

	// RHS must be a type assertion
	typeAssert, ok := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return ""
	}

	// The type assertion must not be a type switch (x.(type))
	if typeAssert.Type == nil {
		return ""
	}

	// Extract the variable being type-asserted
	ident, ok := typeAssert.X.(*ast.Ident)
	if !ok {
		return ""
	}

	// The condition must be the ok variable (second LHS)
	condIdent, ok := ifStmt.Cond.(*ast.Ident)
	if !ok {
		return ""
	}

	okIdent, ok := assignStmt.Lhs[1].(*ast.Ident)
	if !ok {
		return ""
	}

	if condIdent.Name != okIdent.Name {
		return ""
	}

	return ident.Name
}
