package no_net_lookup_port

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNetLookupPortRule disallows calling net.LookupPort without a context.
type NoNetLookupPortRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupPortRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupPort{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupPortRule) Name() string {
	return "noNetLookupPort"
}

// Group returns the rule group.
func (*NoNetLookupPortRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupPortRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupPort struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupPort) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupPort") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupPort does not accept a context; use (*net.Resolver).LookupPort instead",
	})

	return w
}
