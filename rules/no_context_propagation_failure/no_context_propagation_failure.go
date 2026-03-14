package no_context_propagation_failure

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoContextPropagationFailureRule detects context propagation failures
// that can lead to goroutine or resource leaks.
type NoContextPropagationFailureRule struct{}

// Apply applies the rule to given file.
func (r *NoContextPropagationFailureRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		ctxNames := extractContextParamNames(funcDecl.Type.Params)
		if len(ctxNames) == 0 {
			continue
		}

		w := &lintContextPropagation{
			ctxNames: ctxNames,
			onFailure: func(f lint.Failure) {
				failures = append(failures, f)
			},
		}

		ast.Walk(w, funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoContextPropagationFailureRule) Name() string {
	return "noContextPropagationFailure"
}

// Group returns the rule group.
func (*NoContextPropagationFailureRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoContextPropagationFailureRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// extractContextParamNames returns the names of all context.Context parameters.
func extractContextParamNames(params *ast.FieldList) map[string]bool {
	if params == nil {
		return nil
	}

	names := map[string]bool{}
	for _, field := range params.List {
		if isContextType(field.Type) {
			for _, name := range field.Names {
				names[name.Name] = true
			}
		}
	}
	return names
}

// isContextType checks if an expression represents context.Context.
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "context" && sel.Sel.Name == "Context"
}

type lintContextPropagation struct {
	ctxNames  map[string]bool
	onFailure func(lint.Failure)
}

func (w *lintContextPropagation) Visit(node ast.Node) ast.Visitor {
	goStmt, ok := node.(*ast.GoStmt)
	if !ok {
		return w
	}

	funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
	if !ok {
		// go someFunc(...) — check if context is passed as argument
		if !w.hasContextInArgs(goStmt.Call.Args) {
			w.onFailure(lint.Failure{
				Confidence: 0.8,
				Node:       goStmt,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "goroutine does not propagate the parent context",
			})
		}
		return nil
	}

	// go func(...) { ... }(...) — check if context is used inside the goroutine
	if !w.goroutineUsesContext(funcLit, goStmt.Call.Args) {
		w.onFailure(lint.Failure{
			Confidence: 0.8,
			Node:       goStmt,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "goroutine does not propagate the parent context",
		})
	}

	return nil // don't walk into the goroutine body again
}

// hasContextInArgs checks if any of the arguments is a context variable.
func (w *lintContextPropagation) hasContextInArgs(args []ast.Expr) bool {
	for _, arg := range args {
		if ident, ok := arg.(*ast.Ident); ok {
			if w.ctxNames[ident.Name] {
				return true
			}
		}
	}
	return false
}

// goroutineUsesContext checks if the goroutine function literal uses the context
// either via a parameter or via closure capture.
func (w *lintContextPropagation) goroutineUsesContext(funcLit *ast.FuncLit, callArgs []ast.Expr) bool {
	// Check if context is passed as an argument to the goroutine
	if w.hasContextInArgs(callArgs) {
		return true
	}

	// Check if context is referenced in the goroutine body (closure capture)
	found := false
	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		if found {
			return false
		}
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if w.ctxNames[ident.Name] {
			found = true
			return false
		}
		return true
	})
	return found
}
