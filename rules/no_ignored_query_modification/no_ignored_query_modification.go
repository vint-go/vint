package no_ignored_query_modification

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIgnoredQueryModificationRule detects calls that modify the url.Values map
// returned by (*net/url.URL).Query() without assigning it first. Since Query()
// returns a copy, such modifications are silently lost.
type NoIgnoredQueryModificationRule struct{}

// Apply applies the rule to given file.
func (r *NoIgnoredQueryModificationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoIgnoredQueryMod{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIgnoredQueryModificationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoIgnoredQueryMod{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoIgnoredQueryModificationRule) Name() string {
	return "noIgnoredQueryModification"
}

// Group returns the rule group.
func (*NoIgnoredQueryModificationRule) Group() string {
	return "suspicious"
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoIgnoredQueryModificationRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoIgnoredQueryModificationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNoIgnoredQueryMod struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// modifyingMethods are the url.Values methods that mutate the map.
var modifyingMethods = map[string]bool{
	"Set": true,
	"Add": true,
	"Del": true,
}

func (w *lintNoIgnoredQueryMod) Visit(node ast.Node) ast.Visitor {
	// Look for an expression statement containing a call like u.Query().Set(...)
	exprStmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return w
	}

	outerCall, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	// The function being called should be a selector, e.g. <expr>.Set
	outerSel, ok := outerCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	methodName := outerSel.Sel.Name
	if !modifyingMethods[methodName] {
		return w
	}

	// The receiver of the outer selector should be a call, e.g. <expr>.Query()
	innerCall, ok := outerSel.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	innerSel, ok := innerCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if innerSel.Sel.Name != "Query" {
		return w
	}

	// Use type checking to confirm the receiver is *url.URL
	recvType := w.pkg.TypeOf(innerSel.X)
	if recvType == nil {
		return w
	}

	if !isURLType(recvType) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       exprStmt,
		Failure:    "(*net/url.URL).Query returns a copy, modifying it doesn't change the URL",
	})

	return w
}

// isURLType checks whether the given type is *url.URL from net/url.
func isURLType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "URL" && obj.Pkg() != nil && obj.Pkg().Path() == "net/url"
}
