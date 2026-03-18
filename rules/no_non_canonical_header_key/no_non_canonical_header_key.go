package no_non_canonical_header_key

import (
	"fmt"
	"go/ast"
	"go/types"
	"net/textproto"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNonCanonicalHeaderKeyRule flags non-canonical keys used in http.Header map index expressions.
type NoNonCanonicalHeaderKeyRule struct{}

// Apply applies the rule to given file.
func (r *NoNonCanonicalHeaderKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNonCanonicalHeaderKey{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNonCanonicalHeaderKeyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNonCanonicalHeaderKey{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoNonCanonicalHeaderKeyRule) Name() string {
	return "noNonCanonicalHeaderKey"
}

// Group returns the rule group.
func (*NoNonCanonicalHeaderKeyRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNonCanonicalHeaderKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoNonCanonicalHeaderKeyRule) RequiresTypecheck() bool {
	return true
}

type lintNonCanonicalHeaderKey struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNonCanonicalHeaderKey) Visit(node ast.Node) ast.Visitor {
	indexExpr, ok := node.(*ast.IndexExpr)
	if !ok {
		return w
	}

	// Check if the indexed expression is of type http.Header
	exprType := w.pkg.TypeOf(indexExpr.X)
	if exprType == nil {
		return w
	}

	if !isHTTPHeaderType(exprType) {
		return w
	}

	// Check if the index is a string literal
	lit, ok := indexExpr.Index.(*ast.BasicLit)
	if !ok {
		return w
	}

	// Remove quotes from the literal value
	key := lit.Value
	if len(key) < 2 {
		return w
	}
	key = key[1 : len(key)-1] // strip quotes

	canonical := textproto.CanonicalMIMEHeaderKey(key)
	if key == canonical {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       indexExpr,
		Failure:    fmt.Sprintf("non-canonical header key %q, use %q instead", key, canonical),
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
