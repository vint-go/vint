package no_inline_sync_once_func

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInlineSyncOnceFuncRule detects inline calls to sync.OnceFunc, sync.OnceValue, and sync.OnceValues.
type NoInlineSyncOnceFuncRule struct{}

// Apply applies the rule to given file.
func (r *NoInlineSyncOnceFuncRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInlineSyncOnceFunc{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInlineSyncOnceFuncRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInlineSyncOnceFunc{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInlineSyncOnceFuncRule) Name() string {
	return "noInlineSyncOnceFunc"
}

// Group returns the rule group.
func (*NoInlineSyncOnceFuncRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInlineSyncOnceFuncRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInlineSyncOnceFunc struct {
	onFailure func(lint.Failure)
}

// isSyncOnceCall checks if the expression is a call to sync.OnceFunc, sync.OnceValue, or sync.OnceValues.
func isSyncOnceCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	return astutils.IsPkgDotName(call.Fun, "sync", "OnceFunc") ||
		astutils.IsPkgDotName(call.Fun, "sync", "OnceValue") ||
		astutils.IsPkgDotName(call.Fun, "sync", "OnceValues")
}

func (w *lintInlineSyncOnceFunc) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call expression where the function being called
	// is itself a sync.OnceFunc/OnceValue/OnceValues call.
	// This detects patterns like: sync.OnceFunc(func() { ... })()
	if isSyncOnceCall(ce.Fun) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "sync.OnceFunc result is called inline and should be stored in a variable instead",
		})
	}

	return w
}
