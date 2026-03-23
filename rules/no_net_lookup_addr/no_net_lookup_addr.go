package no_net_lookup_addr

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetLookupAddrRule disallows calling net.LookupAddr without a context.
type NoNetLookupAddrRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupAddrRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupAddr{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupAddrRule) Name() string {
	return "noNetLookupAddr"
}

// Group returns the rule group.
func (*NoNetLookupAddrRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupAddrRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupAddr struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupAddr) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupAddr") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupAddr does not accept a context; use (*net.Resolver).LookupAddr instead",
	})

	return w
}
