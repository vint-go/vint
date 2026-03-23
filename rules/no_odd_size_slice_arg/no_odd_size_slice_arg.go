package no_odd_size_slice_arg

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoOddSizeSliceArgRule detects odd-sized slice arguments passed to functions
// expecting even-sized slices (e.g., key-value pairs).
type NoOddSizeSliceArgRule struct{}

// evenArgFunc describes a function that expects an even number of variadic arguments
// starting after a fixed number of leading parameters.
type evenArgFunc struct {
	pkgPath  string // full import path
	typeName string // empty for package-level functions, or the type name for methods
	funcName string // function/method name
	fixedArgs int   // number of fixed (non-variadic) arguments before the pairs
}

// evenArgFuncs is the list of functions that expect even-sized variadic arguments.
var evenArgFuncs = []evenArgFunc{
	// strings.NewReplacer(oldnew ...string)
	{pkgPath: "strings", typeName: "", funcName: "NewReplacer", fixedArgs: 0},

	// log/slog package-level functions
	{pkgPath: "log/slog", typeName: "", funcName: "Info", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "", funcName: "Debug", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "", funcName: "Warn", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "", funcName: "Error", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "", funcName: "Log", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "", funcName: "InfoContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "", funcName: "DebugContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "", funcName: "WarnContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "", funcName: "ErrorContext", fixedArgs: 2},

	// log/slog Logger methods
	{pkgPath: "log/slog", typeName: "Logger", funcName: "Info", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "Debug", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "Warn", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "Error", fixedArgs: 1},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "Log", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "InfoContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "DebugContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "WarnContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "ErrorContext", fixedArgs: 2},
	{pkgPath: "log/slog", typeName: "Logger", funcName: "With", fixedArgs: 0},
}

// Apply applies the rule to given file.
func (r *NoOddSizeSliceArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintOddSizeSliceArg{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoOddSizeSliceArgRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintOddSizeSliceArg{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoOddSizeSliceArgRule) Name() string {
	return "noOddSizeSliceArg"
}

// Group returns the rule group.
func (*NoOddSizeSliceArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoOddSizeSliceArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoOddSizeSliceArgRule) RequiresTypecheck() bool {
	return true
}

type lintOddSizeSliceArg struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintOddSizeSliceArg) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Skip calls using the ellipsis operator (e.g., f(args...))
	if call.Ellipsis.IsValid() {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	funcName := sel.Sel.Name

	// Try to match a package-level function call.
	if ident, ok := sel.X.(*ast.Ident); ok {
		if w.pkg.TypesInfo() != nil {
			obj := w.pkg.TypesInfo().ObjectOf(ident)
			if obj != nil {
				if pkgName, ok := obj.(*types.PkgName); ok {
					pkgPath := pkgName.Imported().Path()
					for _, ef := range evenArgFuncs {
						if ef.typeName != "" {
							continue // skip methods
						}
						if ef.pkgPath == pkgPath && ef.funcName == funcName {
							w.checkEvenArgs(call, ef)
							return w
						}
					}
				}
			}
		}
	}

	// Try to match a method call on a known type.
	if w.pkg.TypesInfo() != nil {
		recvType := w.pkg.TypeOf(sel.X)
		if recvType != nil {
			// Dereference pointer if needed.
			if ptr, ok := recvType.(*types.Pointer); ok {
				recvType = ptr.Elem()
			}
			if named, ok := recvType.(*types.Named); ok {
				obj := named.Obj()
				if obj.Pkg() != nil {
					for _, ef := range evenArgFuncs {
						if ef.typeName == "" {
							continue // skip package-level functions
						}
						if ef.pkgPath == obj.Pkg().Path() && ef.typeName == obj.Name() && ef.funcName == funcName {
							w.checkEvenArgs(call, ef)
							return w
						}
					}
				}
			}
		}
	}

	return w
}

// checkEvenArgs checks that the variadic part of the call has an even number of arguments.
func (w *lintOddSizeSliceArg) checkEvenArgs(call *ast.CallExpr, ef evenArgFunc) {
	totalArgs := len(call.Args)
	variadicArgs := totalArgs - ef.fixedArgs
	if variadicArgs <= 0 {
		return
	}

	// For slog functions, skip slog.Attr arguments since they count as a single slot.
	if ef.pkgPath == "log/slog" {
		w.checkSlogArgs(call, ef.fixedArgs)
		return
	}

	if variadicArgs%2 != 0 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryLogic,
			Failure:    fmt.Sprintf("odd number of arguments passed to %s, expected even number (key-value pairs)", ef.funcName),
		})
	}
}

// checkSlogArgs checks slog-style arguments that can contain slog.Attr values
// (which occupy a single slot) mixed with key-value pairs.
func (w *lintOddSizeSliceArg) checkSlogArgs(call *ast.CallExpr, fixedArgs int) {
	args := call.Args[fixedArgs:]
	if len(args) == 0 {
		return
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		// If this argument is a slog.Attr, it takes one slot.
		if w.isSlogAttr(arg) {
			i++
			continue
		}

		// Otherwise this should be a key (string), followed by a value.
		// Check that there is a corresponding value.
		if i+1 >= len(args) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryLogic,
				Failure:    "odd number of arguments passed to slog function, expected even number (key-value pairs)",
			})
			return
		}

		// Consume the key-value pair.
		i += 2
	}
}

// isSlogAttr returns true if the expression has type slog.Attr.
func (w *lintOddSizeSliceArg) isSlogAttr(expr ast.Expr) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj.Name() == "Attr" && obj.Pkg() != nil && obj.Pkg().Path() == "log/slog"
}
