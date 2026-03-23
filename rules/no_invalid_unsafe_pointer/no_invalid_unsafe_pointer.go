package no_invalid_unsafe_pointer

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidUnsafePointerRule checks for invalid conversions of uintptr to unsafe.Pointer.
// The Go specification defines only a few legal patterns for converting between uintptr
// and unsafe.Pointer. Any conversion that does not follow these patterns is invalid and
// may break in future Go versions or on different platforms.
type NoInvalidUnsafePointerRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidUnsafePointerRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkFunction(funcDecl, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoInvalidUnsafePointerRule) Name() string {
	return "noInvalidUnsafePointer"
}

// Group returns the rule group.
func (*NoInvalidUnsafePointerRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidUnsafePointerRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// checkFunction examines a function body for invalid uintptr to unsafe.Pointer conversions.
func checkFunction(funcDecl *ast.FuncDecl, onFailure func(lint.Failure)) {
	// First pass: collect variable names that are assigned from uintptr() calls.
	uintptrVars := collectUintptrVars(funcDecl.Body)

	// Second pass: find unsafe.Pointer(variable) calls where the variable is a stored uintptr.
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !astutils.IsPkgDotName(call.Fun, "unsafe", "Pointer") {
			return true
		}

		if len(call.Args) != 1 {
			return true
		}

		arg := call.Args[0]

		// Check if the argument is a bare identifier that was assigned from uintptr().
		ident, ok := arg.(*ast.Ident)
		if !ok {
			return true
		}

		if uintptrVars[ident.Name] {
			onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryLogic,
				Failure:    "possibly invalid conversion of uintptr to unsafe.Pointer",
			})
		}

		return true
	})
}

// collectUintptrVars scans a block statement for variables assigned from uintptr() calls.
// It returns a set of variable names that hold uintptr values.
func collectUintptrVars(block *ast.BlockStmt) map[string]bool {
	vars := map[string]bool{}

	ast.Inspect(block, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for i, rhs := range stmt.Rhs {
				if isUintptrConversion(rhs) {
					if i < len(stmt.Lhs) {
						if ident, ok := stmt.Lhs[i].(*ast.Ident); ok {
							vars[ident.Name] = true
						}
					}
				}
			}
		case *ast.ValueSpec:
			// var x = uintptr(...)
			for i, val := range stmt.Values {
				if isUintptrConversion(val) {
					if i < len(stmt.Names) {
						vars[stmt.Names[i].Name] = true
					}
				}
			}
		}
		return true
	})

	return vars
}

// isUintptrConversion checks if an expression is a uintptr() type conversion call.
func isUintptrConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "uintptr"
}
