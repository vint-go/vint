package no_http_post

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoHttpPostRule disallows calling http.Post without a context.
type NoHttpPostRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpPostRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttpPost{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoHttpPostRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoHttpPost{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoHttpPostRule) Name() string {
	return "noHttpPost"
}

// Group returns the rule group.
func (*NoHttpPostRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpPostRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttpPost struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttpPost) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "Post") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.Post does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead",
	})

	return w
}
