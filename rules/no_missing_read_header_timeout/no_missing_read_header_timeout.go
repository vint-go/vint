package no_missing_read_header_timeout

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMissingReadHeaderTimeoutRule detects when ReadHeaderTimeout is not configured
// on an http.Server, which can lead to Slowloris-type denial-of-service attacks.
type NoMissingReadHeaderTimeoutRule struct{}

// Apply applies the rule to given file.
func (r *NoMissingReadHeaderTimeoutRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintMissingReadHeaderTimeout{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoMissingReadHeaderTimeoutRule) Name() string {
	return "noMissingReadHeaderTimeout"
}

// Group returns the rule group.
func (*NoMissingReadHeaderTimeoutRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoMissingReadHeaderTimeoutRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMissingReadHeaderTimeout struct {
	onFailure func(lint.Failure)
}

func (w *lintMissingReadHeaderTimeout) Visit(node ast.Node) ast.Visitor {
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !isHTTPServerType(compositeLit.Type) {
		return w
	}

	if hasReadHeaderTimeout(compositeLit.Elts) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       compositeLit,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "http.Server missing ReadHeaderTimeout, which can lead to Slowloris attacks",
	})

	return w
}

// isHTTPServerType checks if the expression refers to http.Server.
// Matches both &http.Server{} and http.Server{}.
func isHTTPServerType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "http" && sel.Sel.Name == "Server"
}

// hasReadHeaderTimeout checks if any of the composite literal elements
// set the ReadHeaderTimeout field.
func hasReadHeaderTimeout(elts []ast.Expr) bool {
	for _, elt := range elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == "ReadHeaderTimeout" {
			return true
		}
	}

	return false
}
