package no_unmarshalable_type

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnmarshalableTypeRule checks for attempts to marshal channels or functions,
// which have no meaningful serialized representation and will cause runtime errors.
type NoUnmarshalableTypeRule struct{}

// marshalFunc describes a function that marshals a value and cannot handle
// channels or functions.
type marshalFunc struct {
	pkg    string // import path
	name   string // function or method name
	argIdx int    // 0-based index of the argument being marshaled
	isFunc bool   // true for package-level func, false for method
}

// knownMarshalFuncs lists the well-known marshal/encode functions and
// which argument index is the value being marshaled.
var knownMarshalFuncs = []marshalFunc{
	{pkg: "encoding/json", name: "Marshal", argIdx: 0, isFunc: true},
	{pkg: "encoding/json", name: "MarshalIndent", argIdx: 0, isFunc: true},
	{pkg: "encoding/xml", name: "Marshal", argIdx: 0, isFunc: true},
	{pkg: "encoding/xml", name: "MarshalIndent", argIdx: 0, isFunc: true},
	{pkg: "encoding/json", name: "Encode", argIdx: 0, isFunc: false},
	{pkg: "encoding/xml", name: "Encode", argIdx: 0, isFunc: false},
	{pkg: "encoding/gob", name: "Encode", argIdx: 0, isFunc: false},
}

// Apply applies the rule to given file.
func (r *NoUnmarshalableTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoUnmarshalableType{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnmarshalableTypeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoUnmarshalableType{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnmarshalableTypeRule) Name() string {
	return "noUnmarshalableType"
}

// Group returns the rule group.
func (*NoUnmarshalableTypeRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnmarshalableTypeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoUnmarshalableTypeRule) RequiresTypecheck() bool {
	return true
}

type lintNoUnmarshalableType struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoUnmarshalableType) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, mf := range knownMarshalFuncs {
		if mf.isFunc {
			if w.matchesPkgFunc(call, mf) {
				w.checkArg(call, mf.argIdx)
				return w
			}
		} else {
			if w.matchesMethod(call, mf) {
				w.checkArg(call, mf.argIdx)
				return w
			}
		}
	}

	return w
}

// matchesPkgFunc checks if the call expression is a call to pkgName.funcName
// using type info to resolve the package path.
func (w *lintNoUnmarshalableType) matchesPkgFunc(call *ast.CallExpr, mf marshalFunc) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != mf.name {
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

	return pkgName.Imported().Path() == mf.pkg
}

// matchesMethod checks if the call expression is a method call where the
// receiver's type belongs to the expected package and the method name matches.
func (w *lintNoUnmarshalableType) matchesMethod(call *ast.CallExpr, mf marshalFunc) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != mf.name {
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

	return fn.Pkg().Path() == mf.pkg
}

// checkArg checks whether the argument at argIdx in the call has a type that
// contains a channel or function, which cannot be marshaled.
func (w *lintNoUnmarshalableType) checkArg(call *ast.CallExpr, argIdx int) {
	if argIdx >= len(call.Args) {
		return
	}

	arg := call.Args[argIdx]
	t := w.pkg.TypeOf(arg)
	if t == nil {
		return
	}

	if containsUnmarshalableType(t) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "cannot marshal channels or functions",
		})
	}
}

// containsUnmarshalableType checks if the type contains a channel or function type.
func containsUnmarshalableType(t types.Type) bool {
	return containsUnmarshalableTypeRecursive(t, make(map[types.Type]bool))
}

func containsUnmarshalableTypeRecursive(t types.Type, visited map[types.Type]bool) bool {
	if visited[t] {
		return false
	}
	visited[t] = true

	switch u := t.Underlying().(type) {
	case *types.Chan:
		return true
	case *types.Signature:
		return true
	case *types.Pointer:
		return containsUnmarshalableTypeRecursive(u.Elem(), visited)
	case *types.Slice:
		return containsUnmarshalableTypeRecursive(u.Elem(), visited)
	case *types.Array:
		return containsUnmarshalableTypeRecursive(u.Elem(), visited)
	case *types.Map:
		return containsUnmarshalableTypeRecursive(u.Key(), visited) ||
			containsUnmarshalableTypeRecursive(u.Elem(), visited)
	case *types.Interface:
		// Interface values can hold anything at runtime; we can't statically determine
		// whether they contain an unmarshalable type, so we don't flag them.
		return false
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if containsUnmarshalableTypeRecursive(u.Field(i).Type(), visited) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
