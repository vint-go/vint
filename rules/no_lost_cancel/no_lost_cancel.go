package no_lost_cancel

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoLostCancelRule checks for failure to call a context cancellation function.
type NoLostCancelRule struct{}

// contextFuncs lists the context package functions that return a cancel function.
var contextFuncs = []string{"WithCancel", "WithTimeout", "WithDeadline"}

// Apply applies the rule to given file.
func (r *NoLostCancelRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		failures = append(failures, checkFunction(funcDecl.Body)...)
	}

	return failures
}

func checkFunction(body *ast.BlockStmt) []lint.Failure {
	var failures []lint.Failure

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Check if the RHS is a context.With* call
		if len(assign.Rhs) != 1 {
			return true
		}
		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isContextWithCall(call) {
			return true
		}

		// The context.With* functions return (ctx, cancel).
		// We expect the LHS to have exactly 2 values.
		if len(assign.Lhs) != 2 {
			return true
		}

		cancelIdent, ok := assign.Lhs[1].(*ast.Ident)
		if !ok {
			return true
		}

		// Case 1: cancel is assigned to blank identifier
		if cancelIdent.Name == "_" {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Node:       assign,
				Category:   lint.FailureCategoryLogic,
				Failure:    "the cancel function returned by context.With* must be called, not discarded",
			})
			return true
		}

		// Case 2: cancel is assigned to a named variable - check if defer cancel() exists
		if !hasDeferCancel(body, cancelIdent.Name) {
			failures = append(failures, lint.Failure{
				Confidence: 0.8,
				Node:       assign,
				Category:   lint.FailureCategoryLogic,
				Failure:    "the cancel function should be deferred immediately to avoid a context leak",
			})
		}

		return true
	})

	return failures
}

func isContextWithCall(call *ast.CallExpr) bool {
	for _, funcName := range contextFuncs {
		if astutils.IsPkgDotName(call.Fun, "context", funcName) {
			return true
		}
	}
	return false
}

func hasDeferCancel(body *ast.BlockStmt, cancelName string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		deferStmt, ok := n.(*ast.DeferStmt)
		if !ok {
			return true
		}
		callExpr := deferStmt.Call
		ident, ok := callExpr.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == cancelName {
			found = true
			return false
		}
		return true
	})
	return found
}

// Name returns the rule name.
func (*NoLostCancelRule) Name() string {
	return "noLostCancel"
}

// Group returns the rule group.
func (*NoLostCancelRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoLostCancelRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
