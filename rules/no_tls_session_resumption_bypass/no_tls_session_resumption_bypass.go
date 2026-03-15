package no_tls_session_resumption_bypass

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoTlsSessionResumptionBypassRule detects TLS session resumption configurations
// that bypass VerifyPeerCertificate when VerifyConnection is not set.
type NoTlsSessionResumptionBypassRule struct{}

// Apply applies the rule to given file.
func (r *NoTlsSessionResumptionBypassRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTlsSessionResumption{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoTlsSessionResumptionBypassRule) Name() string {
	return "noTlsSessionResumptionBypass"
}

// Group returns the rule group.
func (*NoTlsSessionResumptionBypassRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoTlsSessionResumptionBypassRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTlsSessionResumption struct {
	onFailure func(lint.Failure)
}

func (w *lintTlsSessionResumption) Visit(node ast.Node) ast.Visitor {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(cl.Type, "tls", "Config") {
		return w
	}

	hasVerifyPeerCertificate := false
	hasVerifyConnection := false
	hasSessionTicketsDisabled := false

	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch ident.Name {
		case "VerifyPeerCertificate":
			hasVerifyPeerCertificate = true
		case "VerifyConnection":
			hasVerifyConnection = true
		case "SessionTicketsDisabled":
			if astutils.IsIdent(kv.Value, "true") {
				hasSessionTicketsDisabled = true
			}
		}
	}

	if hasVerifyPeerCertificate && !hasVerifyConnection && !hasSessionTicketsDisabled {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "VerifyPeerCertificate is not called during TLS session resumption; set VerifyConnection or disable session tickets",
		})
	}

	return w
}
