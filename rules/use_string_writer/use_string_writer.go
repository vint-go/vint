package use_string_writer

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseStringWriterRule detects Write([]byte("...")) calls that can be replaced
// with WriteString or io.WriteString for better performance.
type UseStringWriterRule struct{}

// Apply applies the rule to the given file.
func (r *UseStringWriterRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseStringWriter{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseStringWriterRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseStringWriter{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseStringWriterRule) Name() string {
	return "useStringWriter"
}

// Group returns the rule group.
func (*UseStringWriterRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseStringWriterRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseStringWriter struct {
	onFailure func(lint.Failure)
}

func (w *lintUseStringWriter) Visit(node ast.Node) ast.Visitor {
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

	// The argument to []byte() should be a string literal
	lit, ok := byteConv.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind.String() != "STRING" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryOptimization,
		Failure:    "use WriteString or io.WriteString instead of Write([]byte(\"...\"))",
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
