package no_slog_key_value_mismatch

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSlogKeyValueMismatchRule checks for mismatched key-value pairs in log/slog calls.
type NoSlogKeyValueMismatchRule struct{}

// slogFuncs maps slog top-level function names to the index of the first
// key-value argument (i.e. the number of fixed leading arguments).
// For example, slog.Info(msg, args...) has 1 fixed arg (the message).
var slogFuncs = map[string]int{
	"Debug":        1,
	"DebugContext": 2,
	"Error":        1,
	"ErrorContext": 2,
	"Info":         1,
	"InfoContext":  2,
	"Log":          2,
	"Warn":         1,
	"WarnContext":  2,
}

// slogMethodFuncs maps slog.Logger method names to the index of the first
// key-value argument.
var slogMethodFuncs = map[string]int{
	"Debug":        1,
	"DebugContext": 2,
	"Error":        1,
	"ErrorContext": 2,
	"Info":         1,
	"InfoContext":  2,
	"Log":          2,
	"Warn":         1,
	"WarnContext":  2,
	"With":         0,
}

// Apply applies the rule to given file.
func (r *NoSlogKeyValueMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintSlogKeyValue{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSlogKeyValueMismatchRule) Name() string {
	return "noSlogKeyValueMismatch"
}

// Group returns the rule group.
func (*NoSlogKeyValueMismatchRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSlogKeyValueMismatchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSlogKeyValueMismatchRule) RequiresTypecheck() bool {
	return true
}

type lintSlogKeyValue struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintSlogKeyValue) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	funcName := sel.Sel.Name

	// Check if this is a call to a slog top-level function.
	if w.isSlogPackageCall(sel) {
		skip, found := slogFuncs[funcName]
		if !found {
			return w
		}
		w.checkArgs(call, skip)
		return w
	}

	// Check if this is a call on a *slog.Logger receiver.
	if w.isSlogLoggerMethod(sel) {
		skip, found := slogMethodFuncs[funcName]
		if !found {
			return w
		}
		w.checkArgs(call, skip)
		return w
	}

	return w
}

// isSlogPackageCall returns true if the selector refers to a function in the
// "log/slog" package (e.g. slog.Info).
func (w *lintSlogKeyValue) isSlogPackageCall(sel *ast.SelectorExpr) bool {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := w.pkg.TypesInfo().ObjectOf(ident)
	if obj == nil {
		return false
	}

	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkgName.Imported().Path() == "log/slog"
}

// isSlogLoggerMethod returns true if the selector is a method call on a
// *slog.Logger value.
func (w *lintSlogKeyValue) isSlogLoggerMethod(sel *ast.SelectorExpr) bool {
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return false
	}

	// Dereference pointer if needed.
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj.Name() == "Logger" && obj.Pkg() != nil && obj.Pkg().Path() == "log/slog"
}

// checkArgs checks the variadic key-value arguments of a slog call.
// skip is the number of leading fixed arguments to ignore.
func (w *lintSlogKeyValue) checkArgs(call *ast.CallExpr, skip int) {
	if call.Ellipsis.IsValid() {
		// If args are expanded with ..., we can't analyze statically.
		return
	}

	args := call.Args[skip:]
	if len(args) == 0 {
		return
	}

	// Walk through arguments, consuming key-value pairs.
	i := 0
	for i < len(args) {
		arg := args[i]

		// If this argument is a slog.Attr, it takes one slot.
		if w.isSlogAttr(arg) {
			i++
			continue
		}

		// Otherwise this should be a key (string), followed by a value.
		// Check that the key is a string.
		if !w.isStringExpr(arg) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       arg,
				Category:   lint.FailureCategoryLogic,
				Failure:    "slog key-value mismatch: non-string key argument",
			})
			return
		}

		// Check that there is a corresponding value.
		if i+1 >= len(args) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       arg,
				Category:   lint.FailureCategoryLogic,
				Failure:    "slog key-value mismatch: missing value for key argument",
			})
			return
		}

		// Consume the key-value pair.
		i += 2
	}
}

// isSlogAttr returns true if the expression has type slog.Attr.
func (w *lintSlogKeyValue) isSlogAttr(expr ast.Expr) bool {
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

// isStringExpr returns true if the expression has an underlying type of string.
func (w *lintSlogKeyValue) isStringExpr(expr ast.Expr) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	basic, ok := t.Underlying().(*types.Basic)
	return ok && basic.Kind() == types.String
}
