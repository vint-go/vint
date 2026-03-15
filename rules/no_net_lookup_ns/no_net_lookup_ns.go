package no_net_lookup_ns

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNetLookupNsRule disallows calling net.LookupNS without a context.
type NoNetLookupNsRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupNsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupNs{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupNsRule) Name() string {
	return "noNetLookupNs"
}

// Group returns the rule group.
func (*NoNetLookupNsRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupNsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupNs struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupNs) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupNS") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupNS does not accept a context; use (*net.Resolver).LookupNS instead",
	})

	return w
}
