package no_integer_overflow_conversion

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIntegerOverflowConversionRule detects type conversions that can lead to integer overflow.
type NoIntegerOverflowConversionRule struct{}

// Apply applies the rule to the given file.
func (r *NoIntegerOverflowConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		boundsCheckedVars := collectBoundsCheckedVars(funcDecl.Body)

		w := &lintIntegerOverflow{
			file:              file,
			boundsCheckedVars: boundsCheckedVars,
			onFailure: func(f lint.Failure) {
				failures = append(failures, f)
			},
		}
		ast.Walk(w, funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoIntegerOverflowConversionRule) Name() string {
	return "noIntegerOverflowConversion"
}

// Group returns the rule group.
func (*NoIntegerOverflowConversionRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoIntegerOverflowConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoIntegerOverflowConversionRule) RequiresTypecheck() bool {
	return true
}

type lintIntegerOverflow struct {
	file              *lint.File
	boundsCheckedVars map[string]bool
	onFailure         func(lint.Failure)
}

func (w *lintIntegerOverflow) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a type conversion (not a function call)
	if len(ce.Args) != 1 {
		return w
	}

	destType := w.file.Pkg.TypeOf(ce.Fun)
	if destType == nil {
		return w
	}

	destBasic, ok := destType.Underlying().(*types.Basic)
	if !ok || destBasic.Info()&types.IsInteger == 0 {
		return w
	}

	srcType := w.file.Pkg.TypeOf(ce.Args[0])
	if srcType == nil {
		return w
	}

	srcBasic, ok := srcType.Underlying().(*types.Basic)
	if !ok || srcBasic.Info()&types.IsInteger == 0 {
		return w
	}

	// Check if the conversion could potentially overflow
	if !canOverflow(srcBasic.Kind(), destBasic.Kind()) {
		return w
	}

	// Check if the argument is a bounds-checked variable
	if ident, ok := ce.Args[0].(*ast.Ident); ok {
		if w.boundsCheckedVars[ident.Name] {
			return w
		}
	}

	srcName := srcBasic.Name()
	destName := destBasic.Name()

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    fmt.Sprintf("potential integer overflow: conversion of %s to %s may cause overflow or sign change", srcName, destName),
	})

	return w
}

// intTypeSize returns the bit size for a basic integer kind.
// For platform-dependent types (int, uint, uintptr), returns 64 as worst case.
func intTypeSize(kind types.BasicKind) int {
	switch kind {
	case types.Int8, types.Uint8:
		return 8
	case types.Int16, types.Uint16:
		return 16
	case types.Int32, types.Uint32:
		return 32
	case types.Int64, types.Uint64:
		return 64
	case types.Int, types.Uint, types.Uintptr:
		return 64 // platform-dependent, assume 64-bit (worst case)
	default:
		return 0
	}
}

// isSigned returns true if the kind is a signed integer type.
func isSigned(kind types.BasicKind) bool {
	switch kind {
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return true
	default:
		return false
	}
}

// canOverflow returns true if converting from src to dest could cause overflow.
func canOverflow(src, dest types.BasicKind) bool {
	srcSize := intTypeSize(src)
	destSize := intTypeSize(dest)

	if srcSize == 0 || destSize == 0 {
		return false
	}

	srcSigned := isSigned(src)
	destSigned := isSigned(dest)

	// Converting to a smaller bit size always risks overflow
	if destSize < srcSize {
		return true
	}

	// Same size but different signedness risks overflow
	if destSize == srcSize && srcSigned != destSigned {
		return true
	}

	// Converting signed to unsigned of same or larger size can overflow with negative values
	if srcSigned && !destSigned {
		return true
	}

	// Converting unsigned to signed of same size can overflow
	if !srcSigned && destSigned && destSize == srcSize {
		return true
	}

	return false
}

// collectBoundsCheckedVars walks a function body and finds variables that
// are compared using relational operators (indicating bounds checking).
func collectBoundsCheckedVars(body *ast.BlockStmt) map[string]bool {
	checked := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		collectComparedVars(ifStmt.Cond, checked)
		return true
	})

	return checked
}

// collectComparedVars extracts variable names from comparison expressions.
func collectComparedVars(expr ast.Expr, checked map[string]bool) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		switch e.Op {
		case token.LSS, token.GTR, token.LEQ, token.GEQ:
			if name := identName(e.X); name != "" {
				checked[name] = true
			}
			if name := identName(e.Y); name != "" {
				checked[name] = true
			}
		case token.LAND, token.LOR:
			collectComparedVars(e.X, checked)
			collectComparedVars(e.Y, checked)
		}
	case *ast.ParenExpr:
		collectComparedVars(e.X, checked)
	}
}

// identName returns the name if expr is an *ast.Ident, "" otherwise.
func identName(expr ast.Expr) string {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}
