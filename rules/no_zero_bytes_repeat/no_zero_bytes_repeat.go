package no_zero_bytes_repeat

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoZeroBytesRepeatRule detects bytes.Repeat calls with a count of 0.
type NoZeroBytesRepeatRule struct{}

// Apply applies the rule to given file.
func (r *NoZeroBytesRepeatRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintZeroBytesRepeat{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoZeroBytesRepeatRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintZeroBytesRepeat{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoZeroBytesRepeatRule) Name() string {
	return "noZeroBytesRepeat"
}

// Group returns the rule group.
func (*NoZeroBytesRepeatRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoZeroBytesRepeatRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintZeroBytesRepeat struct {
	onFailure func(lint.Failure)
}

func (w *lintZeroBytesRepeat) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(call.Fun, "bytes", "Repeat") {
		return w
	}

	if len(call.Args) != 2 {
		return w
	}

	// Check if the second argument is the integer literal 0
	lit, ok := call.Args[1].(*ast.BasicLit)
	if !ok || lit.Kind != token.INT || lit.Value != "0" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryLogic,
		Failure:    "bytes.Repeat called with a count of 0, always returns an empty slice",
	})

	return w
}
