package no_never_nil_check

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNeverNilCheckRule detects comparisons of never-nil values against nil.
// A value that can never be nil (such as a composite literal or the result
// of certain built-in operations) compared to nil is pointless and likely
// indicates a logic error.
type NoNeverNilCheckRule struct{}

// Apply applies the rule to given file.
func (r *NoNeverNilCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	// Walk each function body separately to scope variable tracking.
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		w := &lintNeverNil{
			onFailure: onFailure,
			neverNil:  map[string]bool{},
		}
		walkStmtList(w, fn.Body.List)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNeverNilCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Body == nil {
		return nil
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNeverNil{
		onFailure: onFailure,
		neverNil:  map[string]bool{},
	}
	walkStmtList(w, fn.Body.List)
	return failures
}

// Name returns the rule name.
func (*NoNeverNilCheckRule) Name() string {
	return "noNeverNilCheck"
}

// Group returns the rule group.
func (*NoNeverNilCheckRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoNeverNilCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNeverNil struct {
	onFailure func(lint.Failure)
	// neverNil tracks local variable names known to hold never-nil values.
	neverNil map[string]bool
}

func (w *lintNeverNil) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl, *ast.FuncLit:
		// Don't descend into nested functions; they have their own scope.
		return nil
	case *ast.AssignStmt:
		w.checkAssignment(n)
	case *ast.BinaryExpr:
		w.checkBinaryExpr(n)
	}
	return w
}

// walkStmtList walks a list of statements using the given walker.
func walkStmtList(w *lintNeverNil, stmts []ast.Stmt) {
	for _, stmt := range stmts {
		ast.Walk(w, stmt)
	}
}

// checkAssignment tracks variables assigned from never-nil expressions.
func (w *lintNeverNil) checkAssignment(assign *ast.AssignStmt) {
	if assign.Tok != token.DEFINE && assign.Tok != token.ASSIGN {
		return
	}

	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		if i >= len(assign.Rhs) {
			continue
		}
		rhs := assign.Rhs[i]
		if isNeverNilExpr(rhs) {
			w.neverNil[ident.Name] = true
		} else {
			// If variable is reassigned to something that could be nil, remove tracking.
			delete(w.neverNil, ident.Name)
		}
	}
}

// checkBinaryExpr checks if a binary comparison involves a never-nil value and nil.
func (w *lintNeverNil) checkBinaryExpr(binExpr *ast.BinaryExpr) {
	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return
	}

	var neverNilExpr ast.Expr
	if isNilIdent(binExpr.Y) {
		neverNilExpr = binExpr.X
	} else if isNilIdent(binExpr.X) {
		neverNilExpr = binExpr.Y
	} else {
		return
	}

	if neverNilExpr == nil {
		return
	}

	// Check if the non-nil side is a never-nil expression directly.
	if isNeverNilExpr(neverNilExpr) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryLogic,
			Failure:    fmt.Sprintf("checking never-nil value %s against nil", astutils.GoFmt(neverNilExpr)),
		})
		return
	}

	// Check if the non-nil side is a variable known to be never-nil.
	if ident, ok := neverNilExpr.(*ast.Ident); ok {
		if w.neverNil[ident.Name] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       binExpr,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("checking never-nil value %s against nil", ident.Name),
			})
		}
	}
}

// isNeverNilExpr returns true if the expression is guaranteed to produce a non-nil value.
func isNeverNilExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		// Composite literals of maps, slices, structs, and arrays are never nil.
		return true
	case *ast.UnaryExpr:
		// &x is never nil.
		if e.Op == token.AND {
			return true
		}
	case *ast.CallExpr:
		// new(T) is never nil.
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "new" {
			return true
		}
		// make() for maps, slices, channels is never nil.
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "make" {
			return true
		}
	case *ast.ParenExpr:
		return isNeverNilExpr(e.X)
	}
	return false
}

// isNilIdent checks if the expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
