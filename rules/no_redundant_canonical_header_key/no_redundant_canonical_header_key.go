package no_redundant_canonical_header_key

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantCanonicalHeaderKeyRule flags redundant calls to http.CanonicalHeaderKey
// in arguments to methods on http.Header, since those methods already canonicalize keys.
type NoRedundantCanonicalHeaderKeyRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantCanonicalHeaderKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantCanonicalHeaderKey{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantCanonicalHeaderKeyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantCanonicalHeaderKey{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantCanonicalHeaderKeyRule) Name() string {
	return "noRedundantCanonicalHeaderKey"
}

// Group returns the rule group.
func (*NoRedundantCanonicalHeaderKeyRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantCanonicalHeaderKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoRedundantCanonicalHeaderKeyRule) RequiresTypecheck() bool {
	return true
}

type lintRedundantCanonicalHeaderKey struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// headerMethodsWithKeyArg lists http.Header methods that canonicalize the key internally.
// The value is the index of the key argument (0-based).
var headerMethodsWithKeyArg = map[string]int{
	"Get":    0,
	"Set":    0,
	"Add":    0,
	"Del":    0,
	"Values": 0,
}

func (w *lintRedundantCanonicalHeaderKey) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// We need a method call: selector.Method(args...)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	methodName := sel.Sel.Name
	keyArgIdx, isHeaderMethod := headerMethodsWithKeyArg[methodName]
	if !isHeaderMethod {
		return w
	}

	// Check if the receiver is of type http.Header
	recvType := w.pkg.TypeOf(sel.X)
	if recvType == nil {
		return w
	}

	if !isHTTPHeaderType(recvType) {
		return w
	}

	// Check if there are enough arguments
	if len(call.Args) <= keyArgIdx {
		return w
	}

	// Check if the key argument is a call to http.CanonicalHeaderKey
	keyArg := call.Args[keyArgIdx]
	innerCall, ok := keyArg.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(innerCall.Fun, "http", "CanonicalHeaderKey") {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       innerCall,
		Failure:    "redundant call to http.CanonicalHeaderKey in method call on http.Header",
	})

	return w
}

// isHTTPHeaderType checks if the type is net/http.Header.
func isHTTPHeaderType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "net/http" && obj.Name() == "Header"
}
