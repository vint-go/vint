package no_http_get

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoHttpGetRule disallows calling http.Get without a context.
type NoHttpGetRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpGetRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttpGet{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoHttpGetRule) Name() string {
	return "noHttpGet"
}

// Group returns the rule group.
func (*NoHttpGetRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpGetRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttpGet struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttpGet) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "Get") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.Get does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead",
	})

	return w
}
