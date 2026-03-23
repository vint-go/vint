package no_cgo_pointer_violation

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoCgoPointerViolationRule detects violations of the cgo pointer passing rules.
// When calling C functions from Go, Go code may pass a Go pointer to C provided
// the Go memory to which it points does not contain any Go pointers.
type NoCgoPointerViolationRule struct{}

// Apply applies the rule to given file.
func (r *NoCgoPointerViolationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !importsCgo(file.AST) {
		return nil
	}

	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintCgoPointer{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoCgoPointerViolationRule) Name() string {
	return "noCgoPointerViolation"
}

// Group returns the rule group.
func (*NoCgoPointerViolationRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoCgoPointerViolationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoCgoPointerViolationRule) RequiresTypecheck() bool {
	return true
}

// importsCgo returns true if the file imports the "C" pseudo-package.
func importsCgo(f *ast.File) bool {
	for _, imp := range f.Imports {
		if imp.Path != nil && imp.Path.Value == `"C"` {
			return true
		}
	}
	return false
}

type lintCgoPointer struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintCgoPointer) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call to a C function (C.xxx)
	if !isCCall(call) {
		return w
	}

	// Check each argument for unsafe.Pointer(&goVar) pattern
	for _, arg := range call.Args {
		w.checkArg(arg, call)
	}

	return w
}

// isCCall checks if the call expression is a call to a C function (C.xxx(...)).
func isCCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	return ok && ident.Name == "C"
}

// checkArg examines a single argument to a C function call for cgo pointer violations.
func (w *lintCgoPointer) checkArg(arg ast.Expr, call *ast.CallExpr) {
	// Look for unsafe.Pointer(&expr) pattern
	innerExpr := w.extractUnsafePointerArg(arg)
	if innerExpr == nil {
		return
	}

	// The inner expression should be a unary & (address-of) operation
	unary, ok := innerExpr.(*ast.UnaryExpr)
	if !ok || unary.Op.String() != "&" {
		return
	}

	// Check if the pointed-to type contains Go pointers
	if w.typeContainsGoPointers(unary.X) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "possible violation of cgo pointer passing rules: Go pointer passed to C contains nested Go pointers",
		})
	}
}

// extractUnsafePointerArg returns the inner expression from unsafe.Pointer(expr),
// or nil if the expression is not of that form.
func (w *lintCgoPointer) extractUnsafePointerArg(expr ast.Expr) ast.Expr {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil
	}

	// Check for unsafe.Pointer(...)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "unsafe" || sel.Sel.Name != "Pointer" {
		return nil
	}

	if len(call.Args) != 1 {
		return nil
	}

	return call.Args[0]
}

// typeContainsGoPointers checks whether the type of the given expression
// contains Go pointers (strings, slices, maps, channels, functions, interfaces,
// or pointers).
func (w *lintCgoPointer) typeContainsGoPointers(expr ast.Expr) bool {
	// Try type-checked info first
	t := w.pkg.TypeOf(expr)
	if t != nil {
		return typeHasPointers(t)
	}

	// Fall back to AST-based heuristics
	return w.astTypeContainsGoPointers(expr)
}

// typeHasPointers checks if a types.Type contains Go pointers.
func typeHasPointers(t types.Type) bool {
	return containsPointer(t.Underlying())
}

// containsPointer recursively checks if a type contains Go pointers.
func containsPointer(t types.Type) bool {
	switch u := t.(type) {
	case *types.Basic:
		return u.Kind() == types.String || u.Kind() == types.UnsafePointer
	case *types.Pointer:
		return true
	case *types.Slice:
		return true // slices contain a pointer to underlying array
	case *types.Map:
		return true
	case *types.Chan:
		return true
	case *types.Signature:
		return true // function values contain pointers
	case *types.Interface:
		return true
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if containsPointer(u.Field(i).Type().Underlying()) {
				return true
			}
		}
		return false
	case *types.Array:
		return containsPointer(u.Elem().Underlying())
	case *types.Named:
		return containsPointer(u.Underlying())
	default:
		return false
	}
}

// astTypeContainsGoPointers uses AST heuristics to detect if an expression's type
// likely contains Go pointers (when type info is not available).
func (w *lintCgoPointer) astTypeContainsGoPointers(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Can't determine type from identifier alone without type info;
		// conservatively report as a potential violation
		return true
	case *ast.SelectorExpr:
		// Field access - conservatively assume it could contain pointers
		return true
	case *ast.IndexExpr:
		// Array/slice indexing - conservatively assume it could contain pointers
		return true
	case *ast.CompositeLit:
		// Check the type of the composite literal
		return w.compositeLitContainsPointers(e)
	default:
		return true // conservative default
	}
}

// compositeLitContainsPointers checks if a composite literal type contains Go pointers.
func (w *lintCgoPointer) compositeLitContainsPointers(lit *ast.CompositeLit) bool {
	if lit.Type == nil {
		return true // can't determine, be conservative
	}
	return astTypeExprContainsPointers(lit.Type)
}

// astTypeExprContainsPointers checks if an AST type expression represents a type
// that contains Go pointers.
func astTypeExprContainsPointers(typeExpr ast.Expr) bool {
	switch t := typeExpr.(type) {
	case *ast.Ident:
		// Built-in types
		switch t.Name {
		case "string":
			return true
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128",
			"byte", "rune", "bool", "uintptr":
			return false
		default:
			// User-defined types - conservatively assume they might contain pointers
			return true
		}
	case *ast.ArrayType:
		if t.Len == nil {
			return true // slice: contains Go pointers
		}
		return astTypeExprContainsPointers(t.Elt) // fixed array: check element type
	case *ast.StarExpr:
		return true // pointer type
	case *ast.MapType:
		return true
	case *ast.ChanType:
		return true
	case *ast.FuncType:
		return true
	case *ast.InterfaceType:
		return true
	case *ast.StructType:
		for _, field := range t.Fields.List {
			if astTypeExprContainsPointers(field.Type) {
				return true
			}
		}
		return false
	case *ast.SelectorExpr:
		// Qualified type (e.g., pkg.Type) - conservatively assume pointers
		return true
	default:
		return true // conservative default
	}
}
