package no_http_client_head

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttpClientHeadRule disallows calling (*net/http.Client).Head without a context.
// The method (*http.Client).Head does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*http.Client).Do with a request created via http.NewRequestWithContext instead.
type NoHttpClientHeadRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpClientHeadRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoHttpClientHead{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoHttpClientHeadRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoHttpClientHead{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoHttpClientHeadRule) Name() string {
	return "noHttpClientHead"
}

// Group returns the rule group.
func (*NoHttpClientHeadRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpClientHeadRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoHttpClientHeadRule) RequiresTypecheck() bool {
	return true
}

type lintNoHttpClientHead struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoHttpClientHead) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a selector expression call: <expr>.Head(...)
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Head" {
		return w
	}

	// Check if the receiver type is *net/http.Client
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return w
	}

	if !isHTTPClientType(t) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*http.Client).Head does not accept a context; use (*http.Client).Do with http.NewRequestWithContext instead",
	})

	return w
}

// isHTTPClientType checks if the given type is *net/http.Client.
func isHTTPClientType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Client" && obj.Pkg() != nil && obj.Pkg().Path() == "net/http"
}
