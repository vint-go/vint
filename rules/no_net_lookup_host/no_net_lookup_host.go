package no_net_lookup_host

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetLookupHostRule disallows calling net.LookupHost without a context.
type NoNetLookupHostRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupHostRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupHost{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupHostRule) Name() string {
	return "noNetLookupHost"
}

// Group returns the rule group.
func (*NoNetLookupHostRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupHostRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupHost struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupHost) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupHost") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupHost does not accept a context; use (*net.Resolver).LookupHost instead",
	})

	return w
}
