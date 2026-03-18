package no_unmarshalable_struct

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnmarshalableStructRule detects attempts to marshal structs that have no exported
// fields. Such structs will always marshal to an empty object "{}", which is
// usually not the intended behavior. It also flags unexported fields with
// serialization struct tags (json, xml, yaml), as those tags are silently
// ignored by the encoding packages.
type NoUnmarshalableStructRule struct{}

// marshalFunc describes a function that marshals a value.
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

// Apply applies the rule to the given file.
func (r *NoUnmarshalableStructRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoUnmarshalableStruct{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnmarshalableStructRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoUnmarshalableStruct{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnmarshalableStructRule) Name() string {
	return "noUnmarshalableStruct"
}

// Group returns the rule group.
func (*NoUnmarshalableStructRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnmarshalableStructRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoUnmarshalableStructRule) RequiresTypecheck() bool {
	return true
}

type lintNoUnmarshalableStruct struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoUnmarshalableStruct) Visit(node ast.Node) ast.Visitor {
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
func (w *lintNoUnmarshalableStruct) matchesPkgFunc(call *ast.CallExpr, mf marshalFunc) bool {
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
func (w *lintNoUnmarshalableStruct) matchesMethod(call *ast.CallExpr, mf marshalFunc) bool {
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

// checkArg checks whether the argument at argIdx in the call is a struct type
// with no exported fields that would be marshaled.
func (w *lintNoUnmarshalableStruct) checkArg(call *ast.CallExpr, argIdx int) {
	if argIdx >= len(call.Args) {
		return
	}

	arg := call.Args[argIdx]
	t := w.pkg.TypeOf(arg)
	if t == nil {
		return
	}

	// Unwrap pointer types to get to the underlying struct.
	t = derefPointer(t)

	st := getStructType(t)
	if st == nil {
		return
	}

	// Check if the struct implements the json.Marshaler or xml.Marshaler interface.
	// If it does, the marshaling behavior is custom and we should not flag it.
	if implementsMarshaler(t, w.pkg) {
		return
	}

	// Check if the struct has any exported fields.
	if hasExportedField(st) {
		return
	}

	typeName := getTypeName(t)
	if typeName != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("struct %s has no exported fields and will marshal as an empty object", typeName),
		})
	} else {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "struct has no exported fields and will marshal as an empty object",
		})
	}
}

// derefPointer unwraps pointer types recursively.
func derefPointer(t types.Type) types.Type {
	for {
		p, ok := t.Underlying().(*types.Pointer)
		if !ok {
			break
		}
		t = p.Elem()
	}
	return t
}

// getStructType returns the underlying struct type if t is a struct or
// a named type whose underlying type is a struct. Returns nil otherwise.
func getStructType(t types.Type) *types.Struct {
	st, ok := t.Underlying().(*types.Struct)
	if ok {
		return st
	}
	return nil
}

// hasExportedField checks if a struct has at least one exported field.
// It recurses into embedded structs.
func hasExportedField(st *types.Struct) bool {
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Exported() && !f.Embedded() {
			return true
		}
		// Check embedded struct fields.
		if f.Embedded() {
			embeddedType := derefPointer(f.Type())
			if embSt, ok := embeddedType.Underlying().(*types.Struct); ok {
				if hasExportedField(embSt) {
					return true
				}
			}
		}
	}
	return false
}

// getTypeName returns the name of a named type, or empty string for anonymous types.
func getTypeName(t types.Type) string {
	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name()
	}
	return ""
}

// implementsMarshaler checks if the type (or pointer to type) implements
// json.Marshaler or xml.Marshaler interfaces, indicating custom marshal behavior.
func implementsMarshaler(t types.Type, pkg *lint.Package) bool {
	info := pkg.TypesInfo()
	if info == nil {
		return false
	}

	// Check for MarshalJSON() ([]byte, error) method
	if hasMethod(t, "MarshalJSON") || hasMethod(types.NewPointer(t), "MarshalJSON") {
		return true
	}

	// Check for MarshalXML method
	if hasMethod(t, "MarshalXML") || hasMethod(types.NewPointer(t), "MarshalXML") {
		return true
	}

	return false
}

// hasMethod checks if a type has a method with the given name.
func hasMethod(t types.Type, name string) bool {
	mset := types.NewMethodSet(t)
	for i := 0; i < mset.Len(); i++ {
		if mset.At(i).Obj().Name() == name {
			return true
		}
	}
	return false
}
