package no_invalid_sort_slice_arg

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidSortSliceArgRule checks for calls to sort.Slice that do not pass
// a slice type as the first argument.
type NoInvalidSortSliceArgRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidSortSliceArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintInvalidSortSliceArg{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoInvalidSortSliceArgRule) Name() string {
	return "noInvalidSortSliceArg"
}

// Group returns the rule group.
func (*NoInvalidSortSliceArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidSortSliceArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoInvalidSortSliceArgRule) RequiresTypecheck() bool {
	return true
}

// sortSliceFuncs lists the sort package functions that expect a slice as first argument.
var sortSliceFuncs = map[string]bool{
	"Slice":         true,
	"SliceStable":   true,
	"SliceIsSorted": true,
}

type lintInvalidSortSliceArg struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintInvalidSortSliceArg) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	// Check if this is a call to sort.Slice, sort.SliceStable, or sort.SliceIsSorted
	if !sortSliceFuncs[sel.Sel.Name] {
		return w
	}

	// Verify it's actually the sort package via type info
	typesInfo := w.file.Pkg.TypesInfo()
	if typesInfo != nil {
		obj, ok := typesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || obj.Imported().Path() != "sort" {
			return w
		}
	}

	if len(call.Args) < 1 {
		return w
	}

	argType := w.file.Pkg.TypeOf(call.Args[0])
	if argType == nil {
		return w
	}

	// Get the underlying type to check if it's a slice
	underlying := argType.Underlying()
	if _, ok := underlying.(*types.Slice); ok {
		return w // valid: first argument is a slice
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       call,
		Failure:    fmt.Sprintf("sort.%s's first argument must be a slice; found %s", sel.Sel.Name, argType.String()),
	})

	return w
}
