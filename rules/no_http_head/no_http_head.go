package no_http_head

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttpHeadRule disallows calling http.Head without a context.
type NoHttpHeadRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpHeadRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoHttpHead{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoHttpHeadRule) Name() string {
	return "noHttpHead"
}

// Group returns the rule group.
func (*NoHttpHeadRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpHeadRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoHttpHead struct {
	onFailure func(lint.Failure)
}

func (w *lintNoHttpHead) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "Head") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.Head does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead",
	})

	return w
}
