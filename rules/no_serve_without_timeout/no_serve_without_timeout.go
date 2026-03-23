package no_serve_without_timeout

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoServeWithoutTimeoutRule detects uses of net/http serve functions
// that have no support for setting timeouts.
type NoServeWithoutTimeoutRule struct{}

// Apply applies the rule to given file.
func (r *NoServeWithoutTimeoutRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintServeWithoutTimeout{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoServeWithoutTimeoutRule) Name() string {
	return "noServeWithoutTimeout"
}

// Group returns the rule group.
func (*NoServeWithoutTimeoutRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoServeWithoutTimeoutRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintServeWithoutTimeout struct {
	onFailure func(lint.Failure)
}

func (w *lintServeWithoutTimeout) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "http", "ListenAndServe") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of http.ListenAndServe with no support for setting timeouts, use http.Server with timeouts instead",
		})
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "http", "ListenAndServeTLS") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of http.ListenAndServeTLS with no support for setting timeouts, use http.Server with timeouts instead",
		})
		return w
	}

	return w
}
