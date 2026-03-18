package no_defer_close_before_err_check

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDeferCloseBeforeErrCheckRule warns when Close is deferred before
// checking the error from the assignment that produced the value.
type NoDeferCloseBeforeErrCheckRule struct{}

// Apply applies the rule to given file.
func (r *NoDeferCloseBeforeErrCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeferClose{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDeferCloseBeforeErrCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDeferClose{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDeferCloseBeforeErrCheckRule) Name() string {
	return "noDeferCloseBeforeErrCheck"
}

// Group returns the rule group.
func (*NoDeferCloseBeforeErrCheckRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeferCloseBeforeErrCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDeferClose struct {
	onFailure func(lint.Failure)
}

func (w *lintDeferClose) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil {
			w.checkBlock(n.Body)
		}
		return nil
	case *ast.FuncLit:
		if n.Body != nil {
			w.checkBlock(n.Body)
		}
		return nil
	}
	return w
}

// checkBlock scans a block statement for the pattern:
//
//	x, err := someCall(...)
//	defer x.Close()
//	if err != nil { ... }
//
// The defer Close before the error check is the problem.
func (w *lintDeferClose) checkBlock(block *ast.BlockStmt) {
	// Recursively check nested blocks and function literals first.
	for _, stmt := range block.List {
		ast.Inspect(stmt, func(n ast.Node) bool {
			if fl, ok := n.(*ast.FuncLit); ok {
				if fl.Body != nil {
					w.checkBlock(fl.Body)
				}
				return false
			}
			if bs, ok := n.(*ast.BlockStmt); ok && bs != block {
				w.checkBlock(bs)
				return false
			}
			return true
		})
	}

	stmts := block.List
	for i := 0; i < len(stmts)-1; i++ {
		assign, ok := stmts[i].(*ast.AssignStmt)
		if !ok {
			continue
		}

		// We need a multi-return assignment like `x, err := ...`
		// where one of the LHS variables is named "err" (or "_")
		// and the other is a named variable that gets .Close() deferred.
		if len(assign.Lhs) < 2 {
			continue
		}

		// Find the error variable and resource variables in the LHS.
		hasErr := false
		var resourceNames []string
		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == "err" {
				hasErr = true
			} else if ident.Name != "_" {
				resourceNames = append(resourceNames, ident.Name)
			}
		}

		if !hasErr || len(resourceNames) == 0 {
			continue
		}

		// Check if the next statement is a defer x.Close()
		deferStmt, ok := stmts[i+1].(*ast.DeferStmt)
		if !ok {
			continue
		}

		if isDeferCloseOn(deferStmt, resourceNames) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       deferStmt,
				Category:   lint.FailureCategoryLogic,
				Failure:    "possible nil dereference: Close deferred before checking error from assignment",
			})
		}
	}
}

// isDeferCloseOn checks whether a defer statement calls .Close() on one of the
// given variable names, either directly (defer x.Close()) or via a function
// literal that calls x.Close().
func isDeferCloseOn(deferStmt *ast.DeferStmt, names []string) bool {
	call := deferStmt.Call

	// Direct: defer x.Close()
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok && sel.Sel.Name == "Close" {
		ident, ok := sel.X.(*ast.Ident)
		if ok {
			for _, name := range names {
				if ident.Name == name {
					return true
				}
			}
		}
	}

	// Wrapped: defer func() { x.Close() }() or defer func() { _ = x.Close() }()
	funcLit, ok := call.Fun.(*ast.FuncLit)
	if ok && funcLit.Body != nil {
		for _, stmt := range funcLit.Body.List {
			if containsCloseCall(stmt, names) {
				return true
			}
		}
	}

	return false
}

// containsCloseCall checks whether a statement contains a call to x.Close()
// for any x in names.
func containsCloseCall(stmt ast.Stmt, names []string) bool {
	found := false
	ast.Inspect(stmt, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Close" {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		for _, name := range names {
			if ident.Name == name {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
