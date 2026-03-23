package no_tls_conn_handshake

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTlsConnHandshakeRule disallows calling (*tls.Conn).Handshake without a context.
// Use (*tls.Conn).HandshakeContext instead, which accepts a context.Context.
type NoTlsConnHandshakeRule struct{}

// Apply applies the rule to given file.
func (r *NoTlsConnHandshakeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTlsConnHandshake{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTlsConnHandshakeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTlsConnHandshake{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoTlsConnHandshakeRule) Name() string {
	return "noTlsConnHandshake"
}

// Group returns the rule group.
func (*NoTlsConnHandshakeRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTlsConnHandshakeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoTlsConnHandshakeRule) RequiresTypecheck() bool {
	return true
}

type lintTlsConnHandshake struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintTlsConnHandshake) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Handshake" {
		return w
	}

	// Check if the receiver type is *crypto/tls.Conn
	recvType := w.pkg.TypeOf(sel.X)
	if recvType == nil {
		return w
	}

	if !isTlsConnType(recvType) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*tls.Conn).Handshake does not accept a context; use (*tls.Conn).HandshakeContext instead",
	})

	return w
}

// isTlsConnType checks if the given type is *crypto/tls.Conn.
func isTlsConnType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Conn" && obj.Pkg() != nil && obj.Pkg().Path() == "crypto/tls"
}
