package no_net_lookup_cname

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetLookupCnameRule disallows calling net.LookupCNAME without a context.
type NoNetLookupCnameRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupCnameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupCname{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupCnameRule) Name() string {
	return "noNetLookupCname"
}

// Group returns the rule group.
func (*NoNetLookupCnameRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupCnameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupCname struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupCname) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupCNAME") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupCNAME does not accept a context; use (*net.Resolver).LookupCNAME instead",
	})

	return w
}
