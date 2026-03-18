package no_invalid_binary_arg

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInvalidBinaryArgRule checks that values passed to binary.Read, binary.Write,
// and binary.Size are of types supported by the encoding/binary package.
type NoInvalidBinaryArgRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidBinaryArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoInvalidBinaryArg{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidBinaryArgRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoInvalidBinaryArg{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidBinaryArgRule) Name() string {
	return "noInvalidBinaryArg"
}

// Group returns the rule group.
func (*NoInvalidBinaryArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidBinaryArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*NoInvalidBinaryArgRule) RequiresTypecheck() bool {
	return true
}

// binaryFuncDataArgIndex maps encoding/binary function names to the index of
// the data argument that must be a fixed-size type.
var binaryFuncDataArgIndex = map[string]int{
	"Read":  2, // binary.Read(r, order, data)
	"Write": 2, // binary.Write(w, order, data)
	"Size":  0, // binary.Size(v)
}

type lintNoInvalidBinaryArg struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoInvalidBinaryArg) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	funcName := sel.Sel.Name
	argIdx, known := binaryFuncDataArgIndex[funcName]
	if !known {
		return w
	}

	// Check if this is the encoding/binary package
	if !astutils.IsPkgDotName(call.Fun, "binary", funcName) {
		return w
	}

	// Verify it's actually the encoding/binary package via type info
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	typesInfo := w.file.Pkg.TypesInfo()
	if typesInfo != nil {
		obj, ok := typesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || obj.Imported().Path() != "encoding/binary" {
			return w
		}
	}

	if len(call.Args) <= argIdx {
		return w
	}

	dataArg := call.Args[argIdx]
	argType := w.file.Pkg.TypeOf(dataArg)
	if argType == nil {
		return w
	}

	if !isValidBinaryType(argType) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("unsupported type %s for encoding/binary, must be a fixed-size type", argType),
		})
	}

	return w
}

// isValidBinaryType checks whether the given type is supported by encoding/binary.
// The encoding/binary package supports:
// - bool
// - fixed-size integer types: int8, int16, int32, int64, uint8, uint16, uint32, uint64
// - float32, float64
// - complex64, complex128
// - arrays of valid types
// - structs where all fields are valid types
// - slices of valid types
// - pointers to any of the above
func isValidBinaryType(t types.Type) bool {
	return isValidBinaryTypeInner(t, true)
}

func isValidBinaryTypeInner(t types.Type, allowPointer bool) bool {
	// Dereference pointers
	if ptr, ok := t.(*types.Pointer); ok {
		if !allowPointer {
			return false
		}
		return isValidBinaryTypeInner(ptr.Elem(), false)
	}

	switch u := t.Underlying().(type) {
	case *types.Basic:
		return isFixedSizeBasic(u)
	case *types.Array:
		return isValidBinaryTypeInner(u.Elem(), false)
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if !isValidBinaryTypeInner(u.Field(i).Type(), false) {
				return false
			}
		}
		return true
	case *types.Slice:
		return isValidBinaryTypeInner(u.Elem(), false)
	default:
		return false
	}
}

// isFixedSizeBasic returns true if the basic type is a fixed-size type
// supported by encoding/binary.
func isFixedSizeBasic(b *types.Basic) bool {
	switch b.Kind() {
	case types.Bool,
		types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.Float32, types.Float64,
		types.Complex64, types.Complex128:
		return true
	default:
		return false
	}
}
