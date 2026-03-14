package no_insecure_cookie

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInsecureCookieRule detects insecure HTTP cookie configurations missing
// Secure, HttpOnly, or SameSite attributes.
type NoInsecureCookieRule struct{}

// Apply applies the rule to given file.
func (r *NoInsecureCookieRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInsecureCookie{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoInsecureCookieRule) Name() string {
	return "noInsecureCookie"
}

// Group returns the rule group.
func (*NoInsecureCookieRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoInsecureCookieRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInsecureCookie struct {
	onFailure func(lint.Failure)
}

func (w *lintInsecureCookie) Visit(node ast.Node) ast.Visitor {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !isHTTPCookieType(cl.Type) {
		return w
	}

	// Check which security-relevant fields are set and their values
	hasSecure := false
	secureIsTrue := false
	hasHttpOnly := false
	httpOnlyIsTrue := false
	hasSameSite := false

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
		case "Secure":
			hasSecure = true
			secureIsTrue = isTrueIdent(kv.Value)
		case "HttpOnly":
			hasHttpOnly = true
			httpOnlyIsTrue = isTrueIdent(kv.Value)
		case "SameSite":
			hasSameSite = true
		}
	}

	if !hasSecure || !secureIsTrue {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "http.Cookie should set Secure to true to prevent transmission over unencrypted connections",
		})
	}

	if !hasHttpOnly || !httpOnlyIsTrue {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "http.Cookie should set HttpOnly to true to prevent JavaScript access",
		})
	}

	if !hasSameSite {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "http.Cookie should set SameSite to prevent CSRF attacks",
		})
	}

	return w
}

// isHTTPCookieType checks if the composite literal type is http.Cookie.
func isHTTPCookieType(expr ast.Expr) bool {
	return astutils.IsPkgDotName(expr, "http", "Cookie")
}

// isTrueIdent checks if the expression is the identifier "true".
func isTrueIdent(expr ast.Expr) bool {
	return astutils.IsIdent(expr, "true")
}
