package no_single_arg_append

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSingleArgAppendRule detects append calls with only a single argument.
type NoSingleArgAppendRule struct{}

// Apply applies the rule to given file.
func (r *NoSingleArgAppendRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSingleArgAppend{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoSingleArgAppendRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSingleArgAppend{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSingleArgAppendRule) Name() string {
	return "noSingleArgAppend"
}

// Group returns the rule group.
func (*NoSingleArgAppendRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSingleArgAppendRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSingleArgAppend struct {
	onFailure func(lint.Failure)
}

func (w *lintSingleArgAppend) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	ident, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != "append" {
		return w
	}

	// append with only one argument (just the slice, no elements) has no effect
	if len(ce.Args) == 1 && !ce.Ellipsis.IsValid() {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "append with only one argument has no effect",
		})
	}

	return w
}
