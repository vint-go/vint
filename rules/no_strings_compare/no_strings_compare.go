package no_strings_compare

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoStringsCompareRule detects usage of strings.Compare which is discouraged by the Go documentation.
type NoStringsCompareRule struct{}

// Apply applies the rule to given file.
func (r *NoStringsCompareRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintStringsCompare{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoStringsCompareRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintStringsCompare{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoStringsCompareRule) Name() string {
	return "noStringsCompare"
}

// Group returns the rule group.
func (*NoStringsCompareRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoStringsCompareRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintStringsCompare struct {
	onFailure func(lint.Failure)
}

func (w *lintStringsCompare) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "strings", "Compare") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryStyle,
		Failure:    "use == or < or > operators instead of strings.Compare",
	})

	return w
}
