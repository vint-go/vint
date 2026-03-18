package use_string_conversion_in_print

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseStringConversionInPrintRule detects []byte arguments passed to fmt print
// functions (Print, Println, Fprint, Fprintln, Sprint, Sprintln) where the
// caller likely wants the string representation. Passing []byte directly prints
// the byte values (e.g. [104 101 108 108 111]) rather than the string.
//
// Source: https://staticcheck.dev/docs/checks/#QF1010
type UseStringConversionInPrintRule struct{}

// Apply applies the rule to the given file.
func (r *UseStringConversionInPrintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()

	w := &lintUseStringConversionInPrint{onFailure: onFailure, typesInfo: typesInfo}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseStringConversionInPrintRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()

	w := &lintUseStringConversionInPrint{onFailure: onFailure, typesInfo: typesInfo}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseStringConversionInPrintRule) Name() string {
	return "useStringConversionInPrint"
}

// Group returns the rule group.
func (*UseStringConversionInPrintRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*UseStringConversionInPrintRule) RequiresTypecheck() bool { return true }

// CacheTier returns the cache tier for this rule.
func (*UseStringConversionInPrintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// fmtPrintFuncs contains the fmt functions that print their arguments without
// a format string. When []byte is passed to these, it prints byte values
// instead of the string representation.
var fmtPrintFuncs = map[string]bool{
	"Print":    true,
	"Println":  true,
	"Fprint":   true,
	"Fprintln": true,
	"Sprint":   true,
	"Sprintln": true,
}

type lintUseStringConversionInPrint struct {
	onFailure func(lint.Failure)
	typesInfo *types.Info
}

func (w *lintUseStringConversionInPrint) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	funcName := sel.Sel.Name
	if !fmtPrintFuncs[funcName] {
		return w
	}

	// Verify it's the "fmt" package
	if !astutils.IsIdent(sel.X, "fmt") {
		return w
	}

	// Check each argument for []byte type
	for _, arg := range call.Args {
		if w.isByteSliceArg(arg) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryStyle,
				Failure:    "convert []byte to string before passing to fmt." + funcName,
			})
			// Only report once per call
			break
		}
	}

	return w
}

// isByteSliceArg checks if the expression is a []byte typed expression.
func (w *lintUseStringConversionInPrint) isByteSliceArg(expr ast.Expr) bool {
	// Check for explicit []byte(...) type conversion
	if call, ok := expr.(*ast.CallExpr); ok {
		if isExplicitByteSliceType(call.Fun) {
			return true
		}
	}

	// Check for composite literal []byte{...}
	if cl, ok := expr.(*ast.CompositeLit); ok {
		if isExplicitByteSliceType(cl.Type) {
			return true
		}
	}

	// Use type info if available to detect variables of type []byte
	if w.typesInfo != nil {
		if t := w.typesInfo.TypeOf(expr); t != nil {
			return t.String() == "[]byte"
		}
	}

	return false
}

// isExplicitByteSliceType checks if expr represents the []byte type expression.
func isExplicitByteSliceType(expr ast.Expr) bool {
	arrayType, ok := expr.(*ast.ArrayType)
	if !ok {
		return false
	}
	// Must be a slice (no length expression), not an array
	if arrayType.Len != nil {
		return false
	}
	ident, ok := arrayType.Elt.(*ast.Ident)
	return ok && ident.Name == "byte"
}
