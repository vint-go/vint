package no_invalid_url_parse

import (
	"go/ast"
	"go/token"
	"net/url"
	"strconv"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidUrlParseRule detects obviously invalid URLs passed to net/url.Parse.
type NoInvalidUrlParseRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidUrlParseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidUrlParse{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidUrlParseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInvalidUrlParse{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidUrlParseRule) Name() string {
	return "noInvalidUrlParse"
}

// Group returns the rule group.
func (*NoInvalidUrlParseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidUrlParseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInvalidUrlParse struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidUrlParse) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(call.Fun, "url", "Parse") {
		return w
	}

	if len(call.Args) < 1 {
		return w
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	rawURL, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	// Use net/url.Parse itself to validate the URL, then also check for
	// backslashes which url.Parse is lenient about but are clearly invalid.
	_, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryLogic,
			Failure:    "invalid URL in url.Parse: " + parseErr.Error(),
		})
		return w
	}

	// Check for backslashes which are not valid in URLs but url.Parse accepts
	for _, ch := range rawURL {
		if ch == '\\' {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryLogic,
				Failure:    "invalid URL in url.Parse: backslash is not allowed in URLs",
			})
			return w
		}
	}

	return w
}
