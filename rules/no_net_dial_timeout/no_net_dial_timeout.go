package no_net_dial_timeout

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetDialTimeoutRule disallows calling net.DialTimeout without a context.
type NoNetDialTimeoutRule struct{}

// Apply applies the rule to given file.
func (r *NoNetDialTimeoutRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetDialTimeout{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetDialTimeoutRule) Name() string {
	return "noNetDialTimeout"
}

// Group returns the rule group.
func (*NoNetDialTimeoutRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetDialTimeoutRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetDialTimeout struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetDialTimeout) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "DialTimeout") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.DialTimeout does not accept a context; use (*net.Dialer).DialContext with Timeout instead",
	})

	return w
}
