package no_redundant_string_byte_conversion

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantStringByteConversionRule detects redundant conversions between
// string and []byte when an equivalent function in the other package exists.
type NoRedundantStringByteConversionRule struct{}

// bytesFuncToStringsFunc maps bytes package functions to their strings equivalents.
// When a bytes.* function is called with []byte(s) arguments, the strings
// equivalent should be used to avoid allocation.
var bytesFuncToStringsFunc = map[string]string{
	"Contains":    "strings.Contains",
	"ContainsAny": "strings.ContainsAny",
	"Count":       "strings.Count",
	"EqualFold":   "strings.EqualFold",
	"HasPrefix":   "strings.HasPrefix",
	"HasSuffix":   "strings.HasSuffix",
	"Index":       "strings.Index",
	"IndexAny":    "strings.IndexAny",
	"IndexByte":   "strings.IndexByte",
	"IndexRune":   "strings.IndexRune",
	"LastIndex":   "strings.LastIndex",
	"LastIndexAny": "strings.LastIndexAny",
	"LastIndexByte": "strings.LastIndexByte",
}

// stringsFuncToByteFunc maps strings package functions to their bytes equivalents.
// When a strings.* function is called with string(b) arguments, the bytes
// equivalent should be used to avoid allocation.
var stringsFuncToByteFunc = map[string]string{
	"Contains":      "bytes.Contains",
	"ContainsAny":   "bytes.ContainsAny",
	"Count":         "bytes.Count",
	"EqualFold":     "bytes.EqualFold",
	"HasPrefix":     "bytes.HasPrefix",
	"HasSuffix":     "bytes.HasSuffix",
	"Index":         "bytes.Index",
	"IndexAny":      "bytes.IndexAny",
	"IndexByte":     "bytes.IndexByte",
	"IndexRune":     "bytes.IndexRune",
	"LastIndex":     "bytes.LastIndex",
	"LastIndexAny":  "bytes.LastIndexAny",
	"LastIndexByte": "bytes.LastIndexByte",
}

// Apply applies the rule to given file.
func (r *NoRedundantStringByteConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantStringByteConversion{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantStringByteConversionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintRedundantStringByteConversion{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantStringByteConversionRule) Name() string {
	return "noRedundantStringByteConversion"
}

// Group returns the rule group.
func (*NoRedundantStringByteConversionRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantStringByteConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantStringByteConversion struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantStringByteConversion) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	funcName := sel.Sel.Name

	switch pkgIdent.Name {
	case "bytes":
		// Check if this bytes function has a strings equivalent
		alternative, exists := bytesFuncToStringsFunc[funcName]
		if !exists {
			return w
		}
		// Check if the first argument is a []byte(s) conversion from string
		if len(ce.Args) > 0 && isByteSliceFromStringConversion(ce.Args[0]) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryOptimization,
				Failure:    fmt.Sprintf("redundant []byte conversion in bytes.%s call, use %s instead", funcName, alternative),
			})
		}
	case "strings":
		// Check if this strings function has a bytes equivalent
		alternative, exists := stringsFuncToByteFunc[funcName]
		if !exists {
			return w
		}
		// Check if the first argument is a string(b) conversion from []byte
		if len(ce.Args) > 0 && isStringFromByteSliceConversion(ce.Args[0]) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryOptimization,
				Failure:    fmt.Sprintf("redundant string conversion in strings.%s call, use %s instead", funcName, alternative),
			})
		}
	}

	return w
}

// isByteSliceFromStringConversion checks if the expression is a []byte(expr) conversion.
func isByteSliceFromStringConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// Check it's a []byte() type conversion
	arrayType, ok := call.Fun.(*ast.ArrayType)
	if !ok {
		return false
	}

	// Check the array has no length (i.e. it's a slice)
	if arrayType.Len != nil {
		return false
	}

	// Check the element type is byte
	return astutils.IsIdent(arrayType.Elt, "byte") && len(call.Args) == 1
}

// isStringFromByteSliceConversion checks if the expression is a string(expr) conversion.
func isStringFromByteSliceConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "string" && len(call.Args) == 1
}
