package no_atoi_overflow

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoAtoiOverflowRule detects potential integer overflow when converting
// strconv.Atoi results to smaller integer types.
type NoAtoiOverflowRule struct{}

// Apply applies the rule to the given file.
func (r *NoAtoiOverflowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		atoiVars := collectAtoiVars(funcDecl.Body)
		if len(atoiVars) == 0 {
			continue
		}

		boundsChecked := collectBoundsCheckedVars(funcDecl.Body, atoiVars)

		findUnsafeConversions(funcDecl.Body, atoiVars, boundsChecked, func(f lint.Failure) {
			failures = append(failures, f)
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoAtoiOverflowRule) Name() string {
	return "noAtoiOverflow"
}

// Group returns the rule group.
func (*NoAtoiOverflowRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoAtoiOverflowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// smallIntTypes lists the integer types that are smaller than int
// and thus may overflow when assigned the result of strconv.Atoi.
var smallIntTypes = map[string]bool{
	"int8":   true,
	"int16":  true,
	"int32":  true,
	"uint":   true,
	"uint8":  true,
	"uint16": true,
	"uint32": true,
}

// collectAtoiVars walks a function body and returns a set of variable names
// that receive the first return value from strconv.Atoi calls.
func collectAtoiVars(body *ast.BlockStmt) map[string]bool {
	vars := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Check if the RHS is a call to strconv.Atoi
		if len(assign.Rhs) != 1 {
			return true
		}
		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}
		if !astutils.IsPkgDotName(call.Fun, "strconv", "Atoi") {
			return true
		}

		// The first LHS variable is the int result
		if len(assign.Lhs) >= 1 {
			if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name != "_" {
				vars[ident.Name] = true
			}
		}

		return true
	})

	return vars
}

// collectBoundsCheckedVars finds Atoi result variables that are compared
// using relational operators (indicating bounds checking).
func collectBoundsCheckedVars(body *ast.BlockStmt, atoiVars map[string]bool) map[string]bool {
	checked := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		binExpr, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		// Check for relational comparison operators
		switch binExpr.Op {
		case token.LSS, token.GTR, token.LEQ, token.GEQ:
			// Check if either side references an Atoi variable
			if name := identName(binExpr.X); name != "" && atoiVars[name] {
				checked[name] = true
			}
			if name := identName(binExpr.Y); name != "" && atoiVars[name] {
				checked[name] = true
			}
		}

		return true
	})

	return checked
}

// identName returns the name if expr is an *ast.Ident, "" otherwise.
func identName(expr ast.Expr) string {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

// findUnsafeConversions walks a function body and reports failures where an
// Atoi result variable is directly converted to a smaller integer type
// without a preceding bounds check.
func findUnsafeConversions(body *ast.BlockStmt, atoiVars, boundsChecked map[string]bool, onFailure func(lint.Failure)) {
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a type conversion to a small int type
		targetType := typeConversionTarget(call)
		if targetType == "" || !smallIntTypes[targetType] {
			return true
		}

		// Check if the argument is an Atoi result variable
		if len(call.Args) != 1 {
			return true
		}
		ident, ok := call.Args[0].(*ast.Ident)
		if !ok || !atoiVars[ident.Name] {
			return true
		}

		// Skip if the variable has been bounds-checked
		if boundsChecked[ident.Name] {
			return true
		}

		onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "potential integer overflow: conversion of strconv.Atoi result to " + targetType + " without bounds check",
		})

		return true
	})
}

// typeConversionTarget returns the target type name if the call expression
// is a type conversion to a known small integer type, or "" otherwise.
func typeConversionTarget(call *ast.CallExpr) string {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}
