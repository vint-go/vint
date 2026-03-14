package no_ssrf_via_variable

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
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

	// Check http.NewRequest(method, url, body) - URL is the second argument
	if astutils.IsPkgDotName(ce.Fun, "http", "NewRequest") {
		if len(ce.Args) > 1 {
			w.checkURLArg(ce, ce.Args[1], "NewRequest")
		}
		return w
	}

	// Check http.NewRequestWithContext(ctx, method, url, body) - URL is the third argument
	if astutils.IsPkgDotName(ce.Fun, "http", "NewRequestWithContext") {
		if len(ce.Args) > 2 {
			w.checkURLArg(ce, ce.Args[2], "NewRequestWithContext")
		}
		return w
	}

	return w
}

// checkURLArg reports a failure if the URL argument is not a string literal.
func (w *lintNoSsrfViaVariable) checkURLArg(call *ast.CallExpr, urlArg ast.Expr, funcName string) {
	if isConstantStringExpr(urlArg) {
		return // hardcoded string literal is safe
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    fmt.Sprintf("potential SSRF: URL passed to http.%s is not a hardcoded constant", funcName),
	})
}

// isConstantStringExpr returns true if the expression is a string literal constant.
func isConstantStringExpr(expr ast.Expr) bool {
	return astutils.IsStringLiteral(expr)
}
