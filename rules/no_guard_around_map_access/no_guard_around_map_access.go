package no_guard_around_map_access

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoGuardAroundMapAccessRule detects unnecessary guard around map access.
// Checking `if _, ok := m[k]; ok { v = m[k] }` can be simplified to `v = m[k]`
// since accessing a missing key returns the zero value.
type NoGuardAroundMapAccessRule struct{}

// Apply applies the rule to given file.
func (r *NoGuardAroundMapAccessRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintGuardAroundMapAccess{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoGuardAroundMapAccessRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintGuardAroundMapAccess{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoGuardAroundMapAccessRule) Name() string {
	return "noGuardAroundMapAccess"
}

// Group returns the rule group.
func (*NoGuardAroundMapAccessRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoGuardAroundMapAccessRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintGuardAroundMapAccess struct {
	onFailure func(lint.Failure)
}

func (w *lintGuardAroundMapAccess) Visit(node ast.Node) ast.Visitor {
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

	// First LHS must be blank identifier (_)
	firstLhs, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok || firstLhs.Name != "_" {
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

	// That statement must be an assignment: v = m[key]
	bodyAssign, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
	if !ok {
		return w
	}

	// The body assignment must be a simple assignment (=) with one LHS and one RHS
	if bodyAssign.Tok != token.ASSIGN || len(bodyAssign.Lhs) != 1 || len(bodyAssign.Rhs) != 1 {
		return w
	}

	// The RHS of the body assignment must be an index expression
	bodyIndex, ok := bodyAssign.Rhs[0].(*ast.IndexExpr)
	if !ok {
		return w
	}

	// The map in the body index must match the map in the init index
	if !exprMatch(bodyIndex.X, indexExpr.X) {
		return w
	}

	// The key in the body index must match the key in the init index
	if !exprMatch(bodyIndex.Index, indexExpr.Index) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ifStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    "unnecessary guard around map access",
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
