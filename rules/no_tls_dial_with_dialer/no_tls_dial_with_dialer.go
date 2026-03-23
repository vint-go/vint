package no_tls_dial_with_dialer

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTlsDialWithDialerRule disallows calling tls.DialWithDialer without a context.
type NoTlsDialWithDialerRule struct{}

// Apply applies the rule to given file.
func (r *NoTlsDialWithDialerRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoTlsDialWithDialer{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTlsDialWithDialerRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoTlsDialWithDialer{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTlsDialWithDialerRule) Name() string {
	return "noTlsDialWithDialer"
}

// Group returns the rule group.
func (*NoTlsDialWithDialerRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTlsDialWithDialerRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoTlsDialWithDialer struct {
	onFailure func(lint.Failure)
}

func (w *lintNoTlsDialWithDialer) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "tls", "DialWithDialer") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "tls.DialWithDialer does not accept a context; use (*tls.Dialer).DialContext with NetDialer instead",
	})

	return w
}
