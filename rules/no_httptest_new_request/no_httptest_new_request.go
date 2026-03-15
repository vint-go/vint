package no_httptest_new_request

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttptestNewRequestRule disallows calling httptest.NewRequest without a context.
type NoHttptestNewRequestRule struct{}

// Apply applies the rule to given file.
func (r *NoHttptestNewRequestRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttptestNewRequest{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoHttptestNewRequestRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoHttptestNewRequest{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoHttptestNewRequestRule) Name() string {
	return "noHttptestNewRequest"
}

// Group returns the rule group.
func (*NoHttptestNewRequestRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttptestNewRequestRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttptestNewRequest struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttptestNewRequest) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "httptest", "NewRequest") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "httptest.NewRequest does not accept a context; use httptest.NewRequestWithContext instead",
	})

	return w
}
