package no_string_index_allocation

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoStringIndexAllocationRule detects strings.Index calls where the first argument
// is a []byte-to-string conversion, which causes an unnecessary allocation.
type NoStringIndexAllocationRule struct{}

// Apply applies the rule to given file.
func (r *NoStringIndexAllocationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintStringIndexAlloc{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoStringIndexAllocationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintStringIndexAlloc{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoStringIndexAllocationRule) Name() string {
	return "noStringIndexAllocation"
}

// Group returns the rule group.
func (*NoStringIndexAllocationRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoStringIndexAllocationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintStringIndexAlloc struct {
	onFailure func(lint.Failure)
}

func (w *lintStringIndexAlloc) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "strings", "Index") {
		return w
	}

	// Check that the first argument is a string(b) conversion from []byte
	if len(ce.Args) < 1 {
		return w
	}

	if isStringFromByteSliceConversion(ce.Args[0]) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryOptimization,
			Failure:    "strings.Index called with a byte-to-string conversion, use bytes.Index instead to avoid allocation",
		})
	}

	return w
}

// isStringFromByteSliceConversion checks if the expression is a string(expr)
// call expression that converts a []byte to string.
func isStringFromByteSliceConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// Check it's a string() type conversion
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "string" && len(call.Args) == 1
}
