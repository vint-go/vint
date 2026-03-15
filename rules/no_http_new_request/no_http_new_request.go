package no_http_new_request

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttpNewRequestRule disallows calling http.NewRequest without a context.
type NoHttpNewRequestRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpNewRequestRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttpNewRequest{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoHttpNewRequestRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoHttpNewRequest{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoHttpNewRequestRule) Name() string {
	return "noHttpNewRequest"
}

// Group returns the rule group.
func (*NoHttpNewRequestRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpNewRequestRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttpNewRequest struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttpNewRequest) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "NewRequest") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.NewRequest does not accept a context; use http.NewRequestWithContext instead",
	})

	return w
}
