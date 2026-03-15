package no_printf_format_mismatch

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoPrintfFormatMismatchRule checks that printf format strings match the provided arguments.
type NoPrintfFormatMismatchRule struct{}

// Apply applies the rule to given file.
func (r *NoPrintfFormatMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintPrintfFormatMismatch{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPrintfFormatMismatchRule) Name() string {
	return "noPrintfFormatMismatch"
}

// Group returns the rule group.
func (*NoPrintfFormatMismatchRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoPrintfFormatMismatchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoPrintfFormatMismatchRule) RequiresTypecheck() bool {
	return true
}

type lintPrintfFormatMismatch struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// printfFunc describes a printf-like function.
type printfFunc struct {
	pkg       string // package name (e.g. "fmt", "log")
	name      string // function name (e.g. "Printf", "Sprintf")
	fmtArgIdx int    // index of the format string argument
	isPrintf  bool   // true if it expects a format string
}

// printfFunctions maps package.Function to its format string argument index.
var printfFunctions = []printfFunc{
	// fmt package - printf-style functions
	{pkg: "fmt", name: "Printf", fmtArgIdx: 0, isPrintf: true},
	{pkg: "fmt", name: "Sprintf", fmtArgIdx: 0, isPrintf: true},
	{pkg: "fmt", name: "Fprintf", fmtArgIdx: 1, isPrintf: true},
	{pkg: "fmt", name: "Errorf", fmtArgIdx: 0, isPrintf: true},

	// log package - printf-style functions
	{pkg: "log", name: "Printf", fmtArgIdx: 0, isPrintf: true},
	{pkg: "log", name: "Fatalf", fmtArgIdx: 0, isPrintf: true},
	{pkg: "log", name: "Panicf", fmtArgIdx: 0, isPrintf: true},
}

// nonPrintfFunctions are functions that do NOT take format strings.
// If they're called with a format-string-like first arg + extra args, that's suspicious.
type nonPrintfFunc struct {
	pkg  string
	name string
}

var nonPrintfFunctions = []nonPrintfFunc{
	{pkg: "fmt", name: "Print"},
	{pkg: "fmt", name: "Println"},
	{pkg: "fmt", name: "Fprint"},
	{pkg: "fmt", name: "Fprintln"},
	{pkg: "fmt", name: "Sprint"},
	{pkg: "fmt", name: "Sprintln"},
	{pkg: "log", name: "Print"},
	{pkg: "log", name: "Println"},
	{pkg: "log", name: "Fatal"},
	{pkg: "log", name: "Fatalln"},
	{pkg: "log", name: "Panic"},
	{pkg: "log", name: "Panicln"},
}

func (w *lintPrintfFormatMismatch) Visit(node ast.Node) ast.Visitor {
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

	pkgName := pkgIdent.Name
	funcName := sel.Sel.Name

	// Verify it's actually a package reference using type info if available.
	if w.pkg.TypesInfo() != nil {
		obj := w.pkg.TypesInfo().Uses[pkgIdent]
		if obj == nil {
			return w
		}
		if _, ok := obj.(*types.PkgName); !ok {
			return w
		}
	}

	// Check printf-style functions
	for _, pf := range printfFunctions {
		if pkgName == pf.pkg && funcName == pf.name {
			w.checkPrintfCall(call, pf)
			return w
		}
	}

	// Check non-printf functions for accidental format string usage
	for _, npf := range nonPrintfFunctions {
		if pkgName == npf.pkg && funcName == npf.name {
			w.checkNonPrintfCall(call, npf)
			return w
		}
	}

	return w
}

// checkPrintfCall verifies a printf-style call.
func (w *lintPrintfFormatMismatch) checkPrintfCall(call *ast.CallExpr, pf printfFunc) {
	// Ensure we have enough args for the format string
	if len(call.Args) <= pf.fmtArgIdx {
		return
	}

	fmtArg := call.Args[pf.fmtArgIdx]
	fmtStr, ok := extractStringLiteral(fmtArg)
	if !ok {
		return // can't analyze non-literal format strings
	}

	verbs, err := parseFormatString(fmtStr)
	if err != nil {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("invalid format string: %s", err),
		})
		return
	}

	// Count actual format verbs (excluding %% which consumes no arguments)
	verbCount := 0
	for _, v := range verbs {
		if v.verb != '%' {
			verbCount += v.argCount()
		}
	}

	// Arguments after the format string
	actualArgs := len(call.Args) - pf.fmtArgIdx - 1
	if call.Ellipsis.IsValid() {
		// Variadic call with ..., can't check argument count
		return
	}

	if verbCount != actualArgs {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("format string expects %d argument(s), but %d provided", verbCount, actualArgs),
		})
		return
	}

	// Check type mismatches for each verb
	argIdx := pf.fmtArgIdx + 1
	for _, v := range verbs {
		if v.verb == '%' {
			continue // %% literal percent
		}
		if v.widthFromArg {
			argIdx++
		}
		if v.precFromArg {
			argIdx++
		}
		if argIdx >= len(call.Args) {
			break
		}
		arg := call.Args[argIdx]
		w.checkVerbTypeMatch(call, v.verb, arg, argIdx-pf.fmtArgIdx)
		argIdx++
	}
}

// checkNonPrintfCall checks if a non-printf function is called with format-string-like args.
func (w *lintPrintfFormatMismatch) checkNonPrintfCall(call *ast.CallExpr, npf nonPrintfFunc) {
	if len(call.Args) < 2 {
		return // needs at least format string + one arg to be suspicious
	}

	firstArg := call.Args[0]
	fmtStr, ok := extractStringLiteral(firstArg)
	if !ok {
		return
	}

	if containsFormatVerb(fmtStr) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("%s.%s call has possible formatting directive; use %s.%sf instead", npf.pkg, npf.name, npf.pkg, npf.name),
		})
	}
}

// formatVerb represents a parsed format verb from a format string.
type formatVerb struct {
	verb         rune
	widthFromArg bool // width specified with *
	precFromArg  bool // precision specified with *
}

// argCount returns the number of arguments consumed by this verb.
func (v formatVerb) argCount() int {
	if v.verb == '%' {
		return 0
	}
	count := 1
	if v.widthFromArg {
		count++
	}
	if v.precFromArg {
		count++
	}
	return count
}

// parseFormatString parses a Go format string and extracts format verbs.
func parseFormatString(format string) ([]formatVerb, error) {
	var verbs []formatVerb
	i := 0
	for i < len(format) {
		r, size := utf8.DecodeRuneInString(format[i:])
		if r != '%' {
			i += size
			continue
		}
		i += size // skip '%'
		if i >= len(format) {
			return nil, fmt.Errorf("format string ends with %%")
		}

		v, consumed, err := parseOneVerb(format[i:])
		if err != nil {
			return nil, err
		}
		verbs = append(verbs, v)
		i += consumed
	}
	return verbs, nil
}

// parseOneVerb parses a single format verb starting after the '%'.
// Returns the verb, the number of bytes consumed, and any error.
func parseOneVerb(s string) (formatVerb, int, error) {
	v := formatVerb{}
	i := 0

	// Skip flags: #, 0, -, +, ' '
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '#' || r == '0' || r == '-' || r == '+' || r == ' ' {
			i += size
		} else {
			break
		}
	}

	if i >= len(s) {
		return v, i, fmt.Errorf("incomplete format verb")
	}

	// Width: either a number, or * (argument)
	r, size := utf8.DecodeRuneInString(s[i:])
	if r == '*' {
		v.widthFromArg = true
		i += size
	} else {
		for i < len(s) {
			r, size = utf8.DecodeRuneInString(s[i:])
			if r >= '0' && r <= '9' {
				i += size
			} else {
				break
			}
		}
	}

	if i >= len(s) {
		return v, i, fmt.Errorf("incomplete format verb")
	}

	// Precision: .number or .*
	r, size = utf8.DecodeRuneInString(s[i:])
	if r == '.' {
		i += size
		if i >= len(s) {
			return v, i, fmt.Errorf("incomplete format verb")
		}
		r, size = utf8.DecodeRuneInString(s[i:])
		if r == '*' {
			v.precFromArg = true
			i += size
		} else {
			for i < len(s) {
				r, size = utf8.DecodeRuneInString(s[i:])
				if r >= '0' && r <= '9' {
					i += size
				} else {
					break
				}
			}
		}
	}

	if i >= len(s) {
		return v, i, fmt.Errorf("incomplete format verb")
	}

	// Verb character
	r, size = utf8.DecodeRuneInString(s[i:])
	i += size
	v.verb = r

	// Check for [n] indexed args - skip them (not commonly used)
	// Validate that this is a known verb
	if !isValidVerb(r) {
		return v, i, fmt.Errorf("unknown format verb '%c'", r)
	}

	return v, i, nil
}

// isValidVerb returns true if the rune is a valid Go format verb.
func isValidVerb(r rune) bool {
	switch r {
	case 'v', 'T', 't', 'b', 'c', 'd', 'o', 'O', 'q', 'x', 'X', 'U',
		'e', 'E', 'f', 'F', 'g', 'G', 's', 'p', 'w',
		'%':
		return true
	}
	return false
}

// extractStringLiteral extracts the string value from a string literal expression.
func extractStringLiteral(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", false
	}
	if lit.Kind.String() != "STRING" {
		return "", false
	}
	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return val, true
}

// containsFormatVerb checks if a string contains printf-style format verbs.
func containsFormatVerb(s string) bool {
	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '%' {
			i += size
			if i >= len(s) {
				return false
			}
			next, _ := utf8.DecodeRuneInString(s[i:])
			if next != '%' && isValidVerb(next) {
				return true
			}
		} else {
			i += size
		}
	}
	return false
}

// checkVerbTypeMatch checks if the argument type matches the format verb.
func (w *lintPrintfFormatMismatch) checkVerbTypeMatch(call *ast.CallExpr, verb rune, arg ast.Expr, argNum int) {
	t := w.pkg.TypeOf(arg)
	if t == nil {
		return // type info not available
	}

	// %v, %T, %p, %w accept anything
	if verb == 'v' || verb == 'T' || verb == 'p' {
		return
	}

	// %w is for errors in fmt.Errorf only
	if verb == 'w' {
		return
	}

	underlying := t.Underlying()

	switch verb {
	case 'd', 'o', 'O', 'x', 'X', 'b', 'c', 'U':
		// Integer verbs
		if !isIntegerType(underlying) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       call,
				Failure:    fmt.Sprintf("format verb %%%c expects an integer type, got %s (argument #%d)", verb, t.String(), argNum),
			})
		}
	case 'e', 'E', 'f', 'F', 'g', 'G':
		// Float verbs
		if !isFloatType(underlying) && !isIntegerType(underlying) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       call,
				Failure:    fmt.Sprintf("format verb %%%c expects a numeric type, got %s (argument #%d)", verb, t.String(), argNum),
			})
		}
	case 's', 'q':
		// String verbs
		if !isStringType(underlying) && !isByteSliceType(underlying) && !isErrorType(t) && !implementsStringer(t) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       call,
				Failure:    fmt.Sprintf("format verb %%%c expects a string type, got %s (argument #%d)", verb, t.String(), argNum),
			})
		}
	case 't':
		// Bool verb
		if !isBoolType(underlying) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       call,
				Failure:    fmt.Sprintf("format verb %%t expects a bool type, got %s (argument #%d)", t.String(), argNum),
			})
		}
	}
}

func isIntegerType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsInteger != 0
}

func isFloatType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&(types.IsFloat|types.IsComplex) != 0
}

func isStringType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsString != 0
}

func isBoolType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsBoolean != 0
}

func isByteSliceType(t types.Type) bool {
	sl, ok := t.(*types.Slice)
	if !ok {
		return false
	}
	basic, ok := sl.Elem().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Byte
}

func isErrorType(t types.Type) bool {
	// Check if the type implements the error interface
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return named.Obj().Id() == "_.error"
}

func implementsStringer(t types.Type) bool {
	// Check if the type has a String() string method
	mset := types.NewMethodSet(t)
	for i := 0; i < mset.Len(); i++ {
		m := mset.At(i)
		if m.Obj().Name() == "String" {
			sig, ok := m.Type().(*types.Signature)
			if ok && sig.Params().Len() == 0 && sig.Results().Len() == 1 {
				if basic, ok := sig.Results().At(0).Type().(*types.Basic); ok && basic.Kind() == types.String {
					return true
				}
			}
		}
	}
	// Also check pointer receiver
	if _, ok := t.(*types.Pointer); !ok {
		return implementsStringer(types.NewPointer(t))
	}
	return false
}

// formatVerbName returns a human-readable name for a format verb.
func formatVerbName(verb rune) string {
	switch verb {
	case 'd':
		return "integer"
	case 's':
		return "string"
	case 'f', 'e', 'g':
		return "float"
	case 'v':
		return "default"
	case 't':
		return "bool"
	case 'p':
		return "pointer"
	default:
		return string(verb)
	}
}

// isPrintfFunc checks if a function name ends with "f" (Printf-style naming convention).
func isPrintfFunc(name string) bool {
	return strings.HasSuffix(name, "f")
}
