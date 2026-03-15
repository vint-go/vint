package no_unsafe_redirect_policy

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnsafeRedirectPolicyRule detects unsafe redirect policies that may propagate
// sensitive headers to different domains when following HTTP redirects.
type NoUnsafeRedirectPolicyRule struct{}

// Apply applies the rule to given file.
func (r *NoUnsafeRedirectPolicyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnsafeRedirectPolicy{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoUnsafeRedirectPolicyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnsafeRedirectPolicy{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnsafeRedirectPolicyRule) Name() string {
	return "noUnsafeRedirectPolicy"
}

// Group returns the rule group.
func (*NoUnsafeRedirectPolicyRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnsafeRedirectPolicyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnsafeRedirectPolicy struct {
	onFailure func(lint.Failure)
}

func (w *lintUnsafeRedirectPolicy) Visit(node ast.Node) ast.Visitor {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !isHTTPClientType(cl.Type) {
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

		if ident.Name != "CheckRedirect" {
			continue
		}

		funcLit, ok := kv.Value.(*ast.FuncLit)
		if !ok {
			continue
		}

		if isUnsafeRedirectFunc(funcLit) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       kv,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "unsafe redirect policy may propagate sensitive headers to different domains",
			})
		}
	}

	return w
}

// isHTTPClientType checks if the expression refers to http.Client.
func isHTTPClientType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "http" && sel.Sel.Name == "Client"
}

// isUnsafeRedirectFunc checks if the function literal copies headers from
// via requests without checking whether the redirect target host differs.
func isUnsafeRedirectFunc(funcLit *ast.FuncLit) bool {
	if funcLit.Body == nil {
		return false
	}

	// Check if the function body contains a range over via headers
	// that copies them to the request without a host check.
	hasHeaderCopy := false
	hasHostCheck := false

	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.RangeStmt:
			// Check if ranging over via[...].Header
			if isViaHeaderAccess(stmt.X) {
				// Check if the body assigns to req.Header
				if bodyAssignsToReqHeader(stmt.Body) {
					hasHeaderCopy = true
				}
			}
		case *ast.IfStmt:
			// Check if there is a host comparison (req.URL.Host != via[...].URL.Host)
			if containsHostComparison(stmt.Cond) {
				hasHostCheck = true
			}
		}
		return true
	})

	return hasHeaderCopy && !hasHostCheck
}

// isViaHeaderAccess checks if expr is accessing via[...].Header.
func isViaHeaderAccess(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != "Header" {
		return false
	}

	// Check if the base is via[something]
	indexExpr, ok := sel.X.(*ast.IndexExpr)
	if !ok {
		return false
	}

	ident, ok := indexExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "via"
}

// bodyAssignsToReqHeader checks if a block statement contains an assignment
// to req.Header[...].
func bodyAssignsToReqHeader(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}

	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for _, lhs := range assign.Lhs {
			if isReqHeaderAccess(lhs) {
				found = true
				return false
			}
		}
		return true
	})

	return found
}

// isReqHeaderAccess checks if the expression accesses req.Header[...].
func isReqHeaderAccess(expr ast.Expr) bool {
	indexExpr, ok := expr.(*ast.IndexExpr)
	if !ok {
		return false
	}

	sel, ok := indexExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != "Header" {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "req"
}

// containsHostComparison checks if the expression contains a comparison
// involving URL.Host fields (typical host check in redirect policies).
func containsHostComparison(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		binExpr, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		// Check if either side references .URL.Host
		if hasURLHostSelector(binExpr.X) || hasURLHostSelector(binExpr.Y) {
			found = true
			return false
		}
		return true
	})
	return found
}

// hasURLHostSelector checks if the expression is of the form something.URL.Host.
func hasURLHostSelector(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != "Host" {
		return false
	}

	innerSel, ok := sel.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return innerSel.Sel.Name == "URL"
}
