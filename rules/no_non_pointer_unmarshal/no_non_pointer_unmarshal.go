package no_non_pointer_unmarshal

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNonPointerUnmarshalRule checks for passing non-pointer or non-interface types
// to unmarshal and decode functions.
type NoNonPointerUnmarshalRule struct{}

// unmarshalFunc describes a function that requires a pointer argument.
type unmarshalFunc struct {
	pkg    string // import path
	name   string // function or method name
	argIdx int    // 0-based index of the argument that must be a pointer
	isFunc bool   // true for package-level func, false for method
}

// knownUnmarshalFuncs lists the well-known unmarshal/decode functions and
// which argument index must be a pointer.
var knownUnmarshalFuncs = []unmarshalFunc{
	{pkg: "encoding/json", name: "Unmarshal", argIdx: 1, isFunc: true},
	{pkg: "encoding/xml", name: "Unmarshal", argIdx: 1, isFunc: true},
	{pkg: "encoding/gob", name: "Decode", argIdx: 0, isFunc: false},
	{pkg: "encoding/json", name: "Decode", argIdx: 0, isFunc: false},
	{pkg: "encoding/xml", name: "Decode", argIdx: 0, isFunc: false},
}

// Apply applies the rule to given file.
func (r *NoNonPointerUnmarshalRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoNonPointerUnmarshal{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNonPointerUnmarshalRule) Name() string {
	return "noNonPointerUnmarshal"
}

// Group returns the rule group.
func (*NoNonPointerUnmarshalRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNonPointerUnmarshalRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoNonPointerUnmarshalRule) RequiresTypecheck() bool {
	return true
}

type lintNoNonPointerUnmarshal struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoNonPointerUnmarshal) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, uf := range knownUnmarshalFuncs {
		if uf.isFunc {
			if w.matchesPkgFunc(call, uf) {
				w.checkArg(call, uf.argIdx)
				return w
			}
		} else {
			if w.matchesMethod(call, uf) {
				w.checkArg(call, uf.argIdx)
				return w
			}
		}
	}

	return w
}

// matchesPkgFunc checks if the call expression is a call to pkgName.funcName
// using type info to resolve the package path.
func (w *lintNoNonPointerUnmarshal) matchesPkgFunc(call *ast.CallExpr, uf unmarshalFunc) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != uf.name {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	info := w.pkg.TypesInfo()
	if info == nil {
		return false
	}

	obj := info.ObjectOf(ident)
	if obj == nil {
		return false
	}

	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkgName.Imported().Path() == uf.pkg
}

// matchesMethod checks if the call expression is a method call where the
// receiver's type belongs to the expected package and the method name matches.
func (w *lintNoNonPointerUnmarshal) matchesMethod(call *ast.CallExpr, uf unmarshalFunc) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != uf.name {
		return false
	}

	info := w.pkg.TypesInfo()
	if info == nil {
		return false
	}

	obj := info.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}

	if fn.Pkg() == nil {
		return false
	}

	return fn.Pkg().Path() == uf.pkg
}

// checkArg checks whether the argument at argIdx in the call is a pointer or interface.
// If not, it reports a failure.
func (w *lintNoNonPointerUnmarshal) checkArg(call *ast.CallExpr, argIdx int) {
	if argIdx >= len(call.Args) {
		return
	}

	arg := call.Args[argIdx]
	t := w.pkg.TypeOf(arg)
	if t == nil {
		return
	}

	// Unwrap the type to check the underlying kind
	if isPointerOrInterface(t) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       call,
		Failure:    "call of unmarshal-like function with non-pointer argument",
	})
}

// isPointerOrInterface returns true if t is a pointer, interface, or an unsafe.Pointer.
func isPointerOrInterface(t types.Type) bool {
	switch t.Underlying().(type) {
	case *types.Pointer:
		return true
	case *types.Interface:
		return true
	}
	return false
}
