package no_insecure_tls_config

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInsecureTlsConfigRule detects insecure TLS configurations in crypto/tls.Config.
type NoInsecureTlsConfigRule struct{}

// Apply applies the rule to given file.
func (r *NoInsecureTlsConfigRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInsecureTls{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoInsecureTlsConfigRule) Name() string {
	return "noInsecureTlsConfig"
}

// Group returns the rule group.
func (*NoInsecureTlsConfigRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoInsecureTlsConfigRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// weakCipherSuites contains cipher suites known to be vulnerable.
var weakCipherSuites = map[string]bool{
	"TLS_RSA_WITH_RC4_128_SHA":                true,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           true,
	"TLS_RSA_WITH_AES_128_CBC_SHA":            true,
	"TLS_RSA_WITH_AES_256_CBC_SHA":            true,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         true,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         true,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         true,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        true,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          true,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     true,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    true,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    true,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      true,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      true,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": true,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   true,
}

// insecureMinVersions contains TLS version constants that are below TLS 1.2.
var insecureMinVersions = map[string]bool{
	"VersionTLS10": true,
	"VersionTLS11": true,
	"VersionSSL30": true,
}

// insecureMaxVersions contains TLS version constants that cap max version too low.
var insecureMaxVersions = map[string]bool{
	"VersionTLS10": true,
	"VersionTLS11": true,
	"VersionSSL30": true,
}

type lintInsecureTls struct {
	onFailure func(lint.Failure)
}

func (w *lintInsecureTls) Visit(node ast.Node) ast.Visitor {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !isTlsConfigType(cl.Type) {
		return w
	}

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
		case "InsecureSkipVerify":
			if astutils.IsIdent(kv.Value, "true") {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       kv,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "TLS InsecureSkipVerify set to true disables certificate verification",
				})
			}
		case "MinVersion":
			if w.isInsecureMinVersion(kv.Value) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       kv,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "TLS MinVersion too low, should be at least tls.VersionTLS12",
				})
			}
		case "MaxVersion":
			if w.isInsecureMaxVersion(kv.Value) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       kv,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "TLS MaxVersion too low, should be at least tls.VersionTLS12",
				})
			}
		case "CipherSuites":
			w.checkCipherSuites(kv.Value)
		}
	}

	return w
}

// isTlsConfigType checks if the composite literal type is tls.Config.
func isTlsConfigType(expr ast.Expr) bool {
	return astutils.IsPkgDotName(expr, "tls", "Config")
}

// isInsecureMinVersion checks if the expression is a TLS version below 1.2.
func (w *lintInsecureTls) isInsecureMinVersion(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !astutils.IsIdent(sel.X, "tls") {
		return false
	}
	return insecureMinVersions[sel.Sel.Name]
}

// isInsecureMaxVersion checks if the expression is a TLS version that caps too low.
func (w *lintInsecureTls) isInsecureMaxVersion(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !astutils.IsIdent(sel.X, "tls") {
		return false
	}
	return insecureMaxVersions[sel.Sel.Name]
}

// checkCipherSuites examines a CipherSuites slice literal for weak cipher suites.
func (w *lintInsecureTls) checkCipherSuites(expr ast.Expr) {
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return
	}
	for _, elt := range cl.Elts {
		sel, ok := elt.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if !astutils.IsIdent(sel.X, "tls") {
			continue
		}
		if weakCipherSuites[sel.Sel.Name] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       sel,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "use of weak TLS cipher suite tls." + sel.Sel.Name,
			})
		}
	}
}
