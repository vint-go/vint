package no_net_ip_bytes_equal

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNetIpBytesEqualRule flags usage of bytes.Equal to compare net.IP values.
type NoNetIpBytesEqualRule struct{}

// Apply applies the rule to given file.
func (r *NoNetIpBytesEqualRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNetIpBytesEqual{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNetIpBytesEqualRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNetIpBytesEqual{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNetIpBytesEqualRule) Name() string {
	return "noNetIpBytesEqual"
}

// Group returns the rule group.
func (*NoNetIpBytesEqualRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule needs type info.
func (*NoNetIpBytesEqualRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoNetIpBytesEqualRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNetIpBytesEqual struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNetIpBytesEqual) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call to bytes.Equal
	if !astutils.IsPkgDotName(callExpr.Fun, "bytes", "Equal") {
		return w
	}

	// bytes.Equal takes exactly 2 arguments
	if len(callExpr.Args) != 2 {
		return w
	}

	// Check if either argument is of type net.IP
	arg0Type := w.file.Pkg.TypeOf(callExpr.Args[0])
	arg1Type := w.file.Pkg.TypeOf(callExpr.Args[1])

	if isNetIP(arg0Type) || isNetIP(arg1Type) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       callExpr,
			Failure:    "use net.IP.Equal to compare net.IP values, not bytes.Equal",
		})
	}

	return w
}

func isNetIP(typ types.Type) bool {
	if typ == nil {
		return false
	}
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "net" && obj.Name() == "IP"
}
