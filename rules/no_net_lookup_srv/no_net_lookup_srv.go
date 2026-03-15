package no_net_lookup_srv

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNetLookupSrvRule disallows calling net.LookupSRV without a context.
type NoNetLookupSrvRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupSrvRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupSrv{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupSrvRule) Name() string {
	return "noNetLookupSrv"
}

// Group returns the rule group.
func (*NoNetLookupSrvRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupSrvRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupSrv struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupSrv) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupSRV") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupSRV does not accept a context; use (*net.Resolver).LookupSRV instead",
	})

	return w
}
