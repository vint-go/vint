package no_net_lookup_txt

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetLookupTxtRule disallows calling net.LookupTXT without a context.
type NoNetLookupTxtRule struct{}

// Apply applies the rule to given file.
func (r *NoNetLookupTxtRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetLookupTxt{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetLookupTxtRule) Name() string {
	return "noNetLookupTxt"
}

// Group returns the rule group.
func (*NoNetLookupTxtRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetLookupTxtRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetLookupTxt struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetLookupTxt) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "LookupTXT") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.LookupTXT does not accept a context; use (*net.Resolver).LookupTXT instead",
	})

	return w
}
