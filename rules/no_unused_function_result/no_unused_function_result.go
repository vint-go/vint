package no_unused_function_result

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedFunctionResultRule checks for unused results of calls to pure functions.
type NoUnusedFunctionResultRule struct{}

// pureFuncEntry represents a package-function pair that is known to be pure.
type pureFuncEntry struct {
	pkg  string
	name string
}

// defaultPureFuncs is the list of pure functions whose results must be used.
var defaultPureFuncs = []pureFuncEntry{
	{pkg: "fmt", name: "Errorf"},
	{pkg: "fmt", name: "Sprintf"},
	{pkg: "fmt", name: "Sprint"},
	{pkg: "errors", name: "New"},
	{pkg: "sort", name: "Reverse"},
}

// Apply applies the rule to given file.
func (r *NoUnusedFunctionResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUnusedFunctionResult{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoUnusedFunctionResultRule) Name() string {
	return "noUnusedFunctionResult"
}

// Group returns the rule group.
func (*NoUnusedFunctionResultRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedFunctionResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnusedFunctionResult struct {
	onFailure func(lint.Failure)
}

func (w *lintUnusedFunctionResult) Visit(node ast.Node) ast.Visitor {
	exprStmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return w
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, entry := range defaultPureFuncs {
		if astutils.IsPkgDotName(callExpr.Fun, entry.pkg, entry.name) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       callExpr,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("result of %s.%s call not used", entry.pkg, entry.name),
			})
			return w
		}
	}

	return w
}
