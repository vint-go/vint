package no_tls_dial

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTlsDialRule disallows calling tls.Dial without a context.
type NoTlsDialRule struct{}

// Apply applies the rule to given file.
func (r *NoTlsDialRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoTlsDial{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTlsDialRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoTlsDial{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTlsDialRule) Name() string {
	return "noTlsDial"
}

// Group returns the rule group.
func (*NoTlsDialRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTlsDialRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoTlsDial struct {
	onFailure func(lint.Failure)
}

func (w *lintNoTlsDial) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "tls", "Dial") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "tls.Dial does not accept a context; use (*tls.Dialer).DialContext instead",
	})

	return w
}
