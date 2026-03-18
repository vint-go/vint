package no_missing_return_after_http_error

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoMissingReturnAfterHttpErrorRule detects http.Error calls without a
// following return statement.
type NoMissingReturnAfterHttpErrorRule struct{}

// Apply applies the rule to given file.
func (r *NoMissingReturnAfterHttpErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintMissingReturnAfterHttpError{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMissingReturnAfterHttpErrorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintMissingReturnAfterHttpError{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMissingReturnAfterHttpErrorRule) Name() string {
	return "noMissingReturnAfterHttpError"
}

// Group returns the rule group.
func (*NoMissingReturnAfterHttpErrorRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMissingReturnAfterHttpErrorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMissingReturnAfterHttpError struct {
	onFailure func(lint.Failure)
}

func (w *lintMissingReturnAfterHttpError) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkBlock(block)
	return w
}

// checkBlock inspects statements in a block, looking for http.Error calls
// that are not followed by a return statement.
func (w *lintMissingReturnAfterHttpError) checkBlock(block *ast.BlockStmt) {
	for i, stmt := range block.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		ce, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		if !astutils.IsPkgDotName(ce.Fun, "http", "Error") {
			continue
		}

		// Check if the next statement is a return (or if this is the last
		// statement in the block, which is also missing a return).
		if i+1 < len(block.List) {
			if _, isReturn := block.List[i+1].(*ast.ReturnStmt); isReturn {
				continue // return follows; no problem
			}
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "http.Error call is not followed by a return statement",
		})
	}
}
