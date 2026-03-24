package no_duplicate_argument

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDuplicateArgumentRule detects suspicious duplicated arguments in function calls.
type NoDuplicateArgumentRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateArgumentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateArgument{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateArgumentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateArgument{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateArgumentRule) Name() string {
	return "noDuplicateArgument"
}

// Group returns the rule group.
func (*NoDuplicateArgumentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateArgumentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// argPair specifies which two argument positions to compare for a given function.
type argPair struct {
	i, j int
}

// whitelistedFuncs maps "pkg.Func" to the argument positions to check.
// Only these functions are checked for duplicate arguments.
var whitelistedFuncs = map[string]argPair{
	// builtins
	"copy": {0, 1},

	// cmp
	"cmp.Compare": {0, 1},

	// math
	"math.Dim": {0, 1},
	"math.Max": {0, 1},
	"math.Min": {0, 1},

	// reflect
	"reflect.Copy":      {0, 1},
	"reflect.DeepEqual": {0, 1},

	// maps
	"maps.Equal": {0, 1},

	// slices
	"slices.Compare": {0, 1},
	"slices.Equal":   {0, 1},

	// strings
	"strings.Contains":    {0, 1},
	"strings.Compare":     {0, 1},
	"strings.EqualFold":   {0, 1},
	"strings.HasPrefix":   {0, 1},
	"strings.HasSuffix":   {0, 1},
	"strings.Index":       {0, 1},
	"strings.LastIndex":   {0, 1},
	"strings.Split":       {0, 1},
	"strings.SplitAfter":  {0, 1},
	"strings.SplitAfterN": {0, 1},
	"strings.SplitN":      {0, 1},
	"strings.Replace":     {1, 2},
	"strings.ReplaceAll":  {1, 2},

	// bytes
	"bytes.Contains":    {0, 1},
	"bytes.Compare":     {0, 1},
	"bytes.Equal":       {0, 1},
	"bytes.EqualFold":   {0, 1},
	"bytes.HasPrefix":   {0, 1},
	"bytes.HasSuffix":   {0, 1},
	"bytes.Index":       {0, 1},
	"bytes.LastIndex":   {0, 1},
	"bytes.Split":       {0, 1},
	"bytes.SplitAfter":  {0, 1},
	"bytes.SplitAfterN": {0, 1},
	"bytes.SplitN":      {0, 1},
	"bytes.Replace":     {1, 2},
	"bytes.ReplaceAll":  {1, 2},

	// types
	"types.Identical":          {0, 1},
	"types.IdenticalIgnoreTags": {0, 1},

	// draw
	"draw.Draw": {0, 2},
}

// whitelistedMethods maps method names to check for receiver == first argument.
var whitelistedMethods = map[string]bool{
	"Equal":   true,
	"Equals":  true,
	"Compare": true,
	"Cmp":     true,
}

type lintNoDuplicateArgument struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDuplicateArgument) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Try method call pattern: x.Method(x) where receiver == arg
	if sel, ok := ce.Fun.(*ast.SelectorExpr); ok {
		if whitelistedMethods[sel.Sel.Name] && len(ce.Args) >= 1 {
			recvStr := astutils.GoFmt(sel.X)
			argStr := astutils.GoFmt(ce.Args[0])
			if recvStr != "" && recvStr == argStr && isPure(sel.X) && isPure(ce.Args[0]) {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       ce,
					Failure:    "suspicious duplicated argument: " + recvStr + " appears more than once in the same call",
				})
				return w
			}
		}
	}

	// Try whitelisted function pattern
	funcName := callName(ce)
	pair, ok := whitelistedFuncs[funcName]
	if !ok {
		return w
	}

	if pair.i >= len(ce.Args) || pair.j >= len(ce.Args) {
		return w
	}

	argI := ce.Args[pair.i]
	argJ := ce.Args[pair.j]

	strI := astutils.GoFmt(argI)
	strJ := astutils.GoFmt(argJ)

	if strI == "" || strI != strJ {
		return w
	}

	// Only flag pure expressions
	if !isPure(argI) || !isPure(argJ) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       ce,
		Failure:    "suspicious duplicated argument: " + strI + " appears more than once in the same call",
	})

	return w
}

// callName extracts the function name from a call expression.
// Returns "pkg.Func" for selector expressions and "func" for identifiers.
func callName(ce *ast.CallExpr) string {
	switch fn := ce.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		if id, ok := fn.X.(*ast.Ident); ok {
			return id.Name + "." + fn.Sel.Name
		}
	}
	return ""
}

// isPure returns true if the expression is side-effect-free.
// Function/method calls are considered impure (except type conversions).
func isPure(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isPure(e.X)
	case *ast.IndexExpr:
		return isPure(e.X) && isPure(e.Index)
	case *ast.StarExpr:
		return isPure(e.X)
	case *ast.ParenExpr:
		return isPure(e.X)
	case *ast.UnaryExpr:
		// Channel receive is impure
		if e.Op.String() == "<-" {
			return false
		}
		return isPure(e.X)
	case *ast.BinaryExpr:
		return isPure(e.X) && isPure(e.Y)
	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			if !isPure(elt) {
				return false
			}
		}
		return true
	case *ast.KeyValueExpr:
		return isPure(e.Value)
	case *ast.FuncLit:
		// A function literal itself is pure (it's not being called)
		return true
	case *ast.SliceExpr:
		if !isPure(e.X) {
			return false
		}
		if e.Low != nil && !isPure(e.Low) {
			return false
		}
		if e.High != nil && !isPure(e.High) {
			return false
		}
		if e.Max != nil && !isPure(e.Max) {
			return false
		}
		return true
	case *ast.CallExpr:
		// Type conversions are pure, actual function calls are not
		return isTypeConversion(e)
	default:
		return false
	}
}

// isTypeConversion returns true if the call expression is a type conversion (e.g. int(x)).
func isTypeConversion(ce *ast.CallExpr) bool {
	if len(ce.Args) != 1 {
		return false
	}
	return isTypeExpr(ce.Fun)
}

// isTypeExpr returns true if the expression represents a type name.
func isTypeExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Built-in types and type names start with lowercase in Go,
		// but we can't reliably distinguish without type info.
		// For safety, we treat all single-ident calls as potentially impure
		// unless they're clearly type names.
		return isBuiltinType(e.Name)
	case *ast.SelectorExpr:
		// pkg.Type - could be a type conversion like pkg.MyType(x)
		// Without type info, we can't be sure, so be conservative
		return false
	case *ast.StarExpr:
		return isTypeExpr(e.X)
	case *ast.ArrayType:
		return true
	case *ast.MapType:
		return true
	case *ast.StructType:
		return true
	case *ast.InterfaceType:
		return true
	case *ast.FuncType:
		return true
	case *ast.ChanType:
		return true
	default:
		return false
	}
}

func isBuiltinType(name string) bool {
	switch name {
	case "bool", "byte", "complex64", "complex128",
		"float32", "float64", "int", "int8", "int16", "int32", "int64",
		"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return true
	}
	return false
}
