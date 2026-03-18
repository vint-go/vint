package no_self_referencing_finalizer

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSelfReferencingFinalizerRule detects when a finalizer function passed to
// runtime.SetFinalizer references the finalized object through a closure,
// which creates a reference cycle preventing garbage collection.
type NoSelfReferencingFinalizerRule struct{}

// Apply applies the rule to given file.
func (r *NoSelfReferencingFinalizerRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSelfReferencingFinalizer{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSelfReferencingFinalizerRule) Name() string {
	return "noSelfReferencingFinalizer"
}

// Group returns the rule group.
func (*NoSelfReferencingFinalizerRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSelfReferencingFinalizerRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoSelfReferencingFinalizer struct {
	onFailure func(lint.Failure)
}

func (w *lintNoSelfReferencingFinalizer) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call to runtime.SetFinalizer
	if !astutils.IsPkgDotName(callExpr.Fun, "runtime", "SetFinalizer") {
		return w
	}

	// Must have exactly 2 arguments
	if len(callExpr.Args) != 2 {
		return w
	}

	// Get the first argument (the object being finalized)
	firstArg := callExpr.Args[0]

	// Get the second argument (the finalizer function)
	funcLit, ok := callExpr.Args[1].(*ast.FuncLit)
	if !ok {
		return w
	}

	// Get the identifier for the first argument
	argIdent := extractIdent(firstArg)
	if argIdent == nil || argIdent.Obj == nil {
		return w
	}

	// Collect parameter names of the function literal so we can skip them
	paramNames := collectParamNames(funcLit)

	// Check if the closure body references the finalized object
	if referencesObject(funcLit.Body, argIdent, paramNames) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       callExpr,
			Failure:    "the finalizer references the finalized object, preventing garbage collection",
		})
	}

	return w
}

// extractIdent extracts the root identifier from an expression.
// For example, for &x it returns x, for x it returns x.
func extractIdent(expr ast.Expr) *ast.Ident {
	switch e := expr.(type) {
	case *ast.Ident:
		return e
	case *ast.UnaryExpr:
		// Handle &x
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident
		}
	}
	return nil
}

// collectParamNames collects all parameter names from a function literal.
func collectParamNames(funcLit *ast.FuncLit) map[string]bool {
	params := map[string]bool{}
	if funcLit.Type != nil && funcLit.Type.Params != nil {
		for _, field := range funcLit.Type.Params.List {
			for _, name := range field.Names {
				params[name.Name] = true
			}
		}
	}
	return params
}

// referencesObject checks if the function body references the given identifier
// (by matching the ast.Object), excluding references through parameter names.
func referencesObject(body *ast.BlockStmt, target *ast.Ident, paramNames map[string]bool) bool {
	if body == nil {
		return false
	}

	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		// Skip if this identifier matches a parameter name
		if paramNames[ident.Name] {
			return true
		}
		// Check if it references the same object
		if ident.Obj != nil && ident.Obj == target.Obj {
			found = true
			return false
		}
		return true
	})

	return found
}
