package no_guard_around_delete

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoGuardAroundDeleteRule detects unnecessary guard around call to delete.
type NoGuardAroundDeleteRule struct{}

// Apply applies the rule to given file.
func (r *NoGuardAroundDeleteRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintGuardAroundDelete{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoGuardAroundDeleteRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintGuardAroundDelete{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoGuardAroundDeleteRule) Name() string {
	return "noGuardAroundDelete"
}

// Group returns the rule group.
func (*NoGuardAroundDeleteRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoGuardAroundDeleteRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintGuardAroundDelete struct {
	onFailure func(lint.Failure)
}

func (w *lintGuardAroundDelete) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Must have no else branch
	if ifStmt.Else != nil {
		return w
	}

	// Must have an init statement (the short variable declaration)
	if ifStmt.Init == nil {
		return w
	}

	// Init must be a short variable declaration: _, ok := m[key]
	assignStmt, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return w
	}

	if assignStmt.Tok != token.DEFINE {
		return w
	}

	// Must have exactly 2 LHS variables and 1 RHS expression
	if len(assignStmt.Lhs) != 2 || len(assignStmt.Rhs) != 1 {
		return w
	}

	// RHS must be an index expression (map lookup)
	indexExpr, ok := assignStmt.Rhs[0].(*ast.IndexExpr)
	if !ok {
		return w
	}

	// The condition must be just the "ok" variable
	condIdent, ok := ifStmt.Cond.(*ast.Ident)
	if !ok {
		return w
	}

	// The condition variable must match the second LHS variable of the assignment
	okIdent, ok := assignStmt.Lhs[1].(*ast.Ident)
	if !ok {
		return w
	}

	if condIdent.Name != okIdent.Name {
		return w
	}

	// Body must contain exactly one statement
	if len(ifStmt.Body.List) != 1 {
		return w
	}

	// That statement must be an expression statement containing a delete call
	exprStmt, ok := ifStmt.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return w
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	// The call must be to "delete"
	funcIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return w
	}

	if funcIdent.Name != "delete" {
		return w
	}

	// delete must have exactly 2 arguments
	if len(callExpr.Args) != 2 {
		return w
	}

	// The map argument in delete must match the map in the index expression
	if !exprMatch(callExpr.Args[0], indexExpr.X) {
		return w
	}

	// The key argument in delete must match the key in the index expression
	if !exprMatch(callExpr.Args[1], indexExpr.Index) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ifStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    "unnecessary guard around call to delete",
	})

	return w
}

// exprMatch checks if two expressions represent the same identifier or selector.
func exprMatch(a, b ast.Expr) bool {
	switch a := a.(type) {
	case *ast.Ident:
		b, ok := b.(*ast.Ident)
		if !ok {
			return false
		}
		return a.Name == b.Name
	case *ast.SelectorExpr:
		b, ok := b.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		return a.Sel.Name == b.Sel.Name && exprMatch(a.X, b.X)
	case *ast.IndexExpr:
		b, ok := b.(*ast.IndexExpr)
		if !ok {
			return false
		}
		return exprMatch(a.X, b.X) && exprMatch(a.Index, b.Index)
	case *ast.BasicLit:
		b, ok := b.(*ast.BasicLit)
		if !ok {
			return false
		}
		return a.Kind == b.Kind && a.Value == b.Value
	case *ast.CallExpr:
		b, ok := b.(*ast.CallExpr)
		if !ok {
			return false
		}
		if !exprMatch(a.Fun, b.Fun) {
			return false
		}
		if len(a.Args) != len(b.Args) {
			return false
		}
		for i := range a.Args {
			if !exprMatch(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
