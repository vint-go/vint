package use_printf_suffix

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UsePrintfSuffixRule checks that printf-like functions are named with 'f' at the end.
type UsePrintfSuffixRule struct{}

// Apply applies the rule to given file.
func (r *UsePrintfSuffixRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if !isPrintfLike(funcDecl) {
			continue
		}

		name := funcDecl.Name.Name
		if strings.HasSuffix(name, "f") {
			continue
		}

		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryNaming,
			Confidence: 1,
			Node:       funcDecl,
			Failure:    fmt.Sprintf("printf-like formatting function '%s' should be named '%sf'", name, name),
		})
	}

	return failures
}

// isPrintfLike checks whether a function declaration matches the printf-like signature:
// 1. No return values
// 2. At least two parameters
// 3. Second-to-last parameter is of type string and named "format"
// 4. Last parameter is variadic with type ...interface{} or ...any
func isPrintfLike(funcDecl *ast.FuncDecl) bool {
	// Must have no return values
	if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
		return false
	}

	params := funcDecl.Type.Params
	if params == nil {
		return false
	}

	// Flatten params to get individual parameter entries
	// Each field in the list can have multiple names
	fields := params.List
	if len(fields) < 2 {
		return false
	}

	// The last field must be variadic (ellipsis)
	lastField := fields[len(fields)-1]
	ellipsis, ok := lastField.Type.(*ast.Ellipsis)
	if !ok {
		return false
	}

	// The variadic type must be interface{} or any
	if !isInterfaceOrAny(ellipsis.Elt) {
		return false
	}

	// The second-to-last field must be of type string and named "format"
	formatField := fields[len(fields)-2]

	// Check the type is "string"
	ident, ok := formatField.Type.(*ast.Ident)
	if !ok || ident.Name != "string" {
		return false
	}

	// Check at least one name is "format"
	if len(formatField.Names) == 0 {
		return false
	}
	hasFormatName := false
	for _, name := range formatField.Names {
		if name.Name == "format" {
			hasFormatName = true
			break
		}
	}
	if !hasFormatName {
		return false
	}

	return true
}

// isInterfaceOrAny checks if the expression is `interface{}` or `any`.
func isInterfaceOrAny(expr ast.Expr) bool {
	// Check for `any` (an identifier)
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "any"
	}

	// Check for `interface{}`
	if iface, ok := expr.(*ast.InterfaceType); ok {
		return iface.Methods == nil || len(iface.Methods.List) == 0
	}

	return false
}

// Name returns the rule name.
func (*UsePrintfSuffixRule) Name() string {
	return "usePrintfSuffix"
}

// Group returns the rule group.
func (*UsePrintfSuffixRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UsePrintfSuffixRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
