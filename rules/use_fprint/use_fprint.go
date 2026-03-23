package use_fprint

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseFprintRule detects fmt.Sprint/fmt.Sprintf calls whose result is
// converted to []byte and passed to an io.Writer's Write method. In
// such cases fmt.Fprint/fmt.Fprintf can be used directly, avoiding the
// intermediate string allocation.
type UseFprintRule struct{}

// Apply applies the rule to the given file.
func (r *UseFprintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseFprint{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseFprintRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseFprint{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseFprintRule) Name() string {
	return "useFprint"
}

// Group returns the rule group.
func (*UseFprintRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseFprintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseFprint struct {
	onFailure func(lint.Failure)
}

func (w *lintUseFprint) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Look for a method call: something.Write(...)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Write" {
		return w
	}

	// Write should have exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	// The argument should be a []byte conversion: []byte(expr)
	byteConv, ok := call.Args[0].(*ast.CallExpr)
	if !ok || len(byteConv.Args) != 1 {
		return w
	}

	if !isByteSliceType(byteConv.Fun) {
		return w
	}

	// The argument to []byte() should be a fmt.Sprint or fmt.Sprintf call
	innerCall, ok := byteConv.Args[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	var sprintFunc string
	if astutils.IsPkgDotName(innerCall.Fun, "fmt", "Sprint") {
		sprintFunc = "fmt.Sprint"
	} else if astutils.IsPkgDotName(innerCall.Fun, "fmt", "Sprintf") {
		sprintFunc = "fmt.Sprintf"
	} else {
		return w
	}

	replacement := "fmt.Fprint"
	if sprintFunc == "fmt.Sprintf" {
		replacement = "fmt.Fprintf"
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryOptimization,
		Failure:    sprintFunc + " result converted to []byte and written; use " + replacement + " instead",
	})

	return w
}

// isByteSliceType checks if expr represents the []byte type expression.
func isByteSliceType(expr ast.Expr) bool {
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
