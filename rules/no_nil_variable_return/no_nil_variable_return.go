package no_nil_variable_return

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNilVariableReturnRule detects return statements that return a variable
// known to be nil from a preceding nil check.
type NoNilVariableReturnRule struct{}

// Apply applies the rule to given file.
func (r *NoNilVariableReturnRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNilVarReturn{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNilVariableReturnRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNilVarReturn{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNilVariableReturnRule) Name() string {
	return "noNilVariableReturn"
}

// Group returns the rule group.
func (*NoNilVariableReturnRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoNilVariableReturnRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNilVarReturn struct {
	onFailure func(lint.Failure)
}

func (w *lintNilVarReturn) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Check for `if x == nil { ... return x ... }` pattern
	w.checkIfStmt(ifStmt)

	return w
}

func (w *lintNilVarReturn) checkIfStmt(ifStmt *ast.IfStmt) {
	// We need a binary expression in the condition: x == nil
	binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return
	}

	// Must be == comparison
	if binExpr.Op != token.EQL {
		return
	}

	// One side must be nil, the other an identifier
	varName := ""
	if isNilIdent(binExpr.Y) {
		if ident, ok := binExpr.X.(*ast.Ident); ok {
			varName = ident.Name
		}
	} else if isNilIdent(binExpr.X) {
		if ident, ok := binExpr.Y.(*ast.Ident); ok {
			varName = ident.Name
		}
	}

	if varName == "" {
		return
	}

	// Check the if body for return statements that return the nil variable
	if ifStmt.Body == nil {
		return
	}

	for _, stmt := range ifStmt.Body.List {
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		for _, result := range retStmt.Results {
			ident, ok := result.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == varName {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       retStmt,
					Failure:    "returning a nil variable '" + varName + "' instead of explicit nil",
				})
			}
		}
	}
}

func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
