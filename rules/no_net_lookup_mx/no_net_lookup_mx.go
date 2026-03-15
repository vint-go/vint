package no_net_lookup_mx

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNetLookupMxRule disallows calling net.LookupMX without a context.
type NoNetLookupMxRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupMxRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupMx{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupMxRule) Name() string {
	return "noNetLookupMx"
}

// Group returns the rule group.
func (*NoNetLookupMxRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupMxRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupMx struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupMx) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupMX") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupMX does not accept a context; use (*net.Resolver).LookupMX instead",
	})

	return w
}
