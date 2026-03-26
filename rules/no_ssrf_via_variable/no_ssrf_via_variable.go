package no_ssrf_via_variable

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSsrfViaVariableRule detects potential SSRF vulnerabilities by identifying
// HTTP requests made with URLs derived from variables or user input.
type NoSsrfViaVariableRule struct{}

// Apply applies the rule to given file.
func (r *NoSsrfViaVariableRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSsrfViaVariable{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSsrfViaVariableRule) Name() string {
	return "noSsrfViaVariable"
}

// Group returns the rule group.
func (*NoSsrfViaVariableRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSsrfViaVariableRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoSsrfViaVariable struct {
	onFailure func(lint.Failure)
}

// httpFuncURLArgIndex maps net/http package-level functions to the
// zero-based index of their URL argument.
var httpFuncURLArgIndex = map[string]int{
	"Get":      0, // http.Get(url)
	"Head":     0, // http.Head(url)
	"Post":     0, // http.Post(url, contentType, body)
	"PostForm": 0, // http.PostForm(url, data)
}

func (w *lintNoSsrfViaVariable) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check package-level http functions: http.Get, http.Head, http.Post, http.PostForm
	for funcName, urlArgIdx := range httpFuncURLArgIndex {
		if astutils.IsPkgDotName(ce.Fun, "http", funcName) {
			if len(ce.Args) > urlArgIdx {
				w.checkURLArg(ce, ce.Args[urlArgIdx], funcName)
			}
			return w
		}
	}

	return w
}

// checkURLArg reports a failure if the URL argument is a variable (non-constant *ast.Ident).
// String literals, constant identifiers, call expressions, and binary expressions are not flagged,
// matching the behavior of gosec G107's ResolveVar.
func (w *lintNoSsrfViaVariable) checkURLArg(call *ast.CallExpr, urlArg ast.Expr, funcName string) {
	// String literals are always safe.
	if astutils.IsStringLiteral(urlArg) {
		return
	}

	// Only flag *ast.Ident nodes that are not constants.
	ident, ok := urlArg.(*ast.Ident)
	if !ok {
		return // call expressions, binary expressions, etc. are not flagged
	}

	// If the identifier refers to a constant, it's safe.
	if ident.Obj != nil && ident.Obj.Kind == ast.Con {
		return
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    fmt.Sprintf("potential SSRF: URL passed to http.%s is not a hardcoded constant", funcName),
	})
}
