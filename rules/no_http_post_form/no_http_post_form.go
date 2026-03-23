package no_http_post_form

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoHttpPostFormRule disallows calling http.PostForm without a context.
type NoHttpPostFormRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpPostFormRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttpPostForm{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoHttpPostFormRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoHttpPostForm{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoHttpPostFormRule) Name() string {
	return "noHttpPostForm"
}

// Group returns the rule group.
func (*NoHttpPostFormRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpPostFormRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttpPostForm struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttpPostForm) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "PostForm") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.PostForm does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead",
	})

	return w
}
