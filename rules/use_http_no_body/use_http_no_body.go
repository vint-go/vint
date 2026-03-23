package use_http_no_body

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseHttpNoBodyRule detects nil usages in http.NewRequest calls and suggests http.NoBody instead.
type UseHttpNoBodyRule struct{}

// Apply applies the rule to given file.
func (r *UseHttpNoBodyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintHttpNoBody{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseHttpNoBodyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintHttpNoBody{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseHttpNoBodyRule) Name() string {
	return "useHttpNoBody"
}

// Group returns the rule group.
func (*UseHttpNoBodyRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseHttpNoBodyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintHttpNoBody struct {
	onFailure func(lint.Failure)
}

func (w *lintHttpNoBody) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "http", "NewRequest") {
		return w
	}

	// http.NewRequest takes 3 arguments: method, url, body
	if len(ce.Args) != 3 {
		return w
	}

	// Check if the third argument (body) is nil
	ident, ok := ce.Args[2].(*ast.Ident)
	if !ok || ident.Name != "nil" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryStyle,
		Failure:    "use http.NoBody instead of nil in http.NewRequest",
	})

	return w
}
