package no_net_lookup_ip

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetLookupIpRule disallows calling net.LookupIP without a context.
type NoNetLookupIpRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupIpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupIp{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupIpRule) Name() string {
	return "noNetLookupIp"
}

// Group returns the rule group.
func (*NoNetLookupIpRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupIpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupIp struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupIp) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupIP") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupIP does not accept a context; use (*net.Resolver).LookupIPAddr instead",
	})

	return w
}
