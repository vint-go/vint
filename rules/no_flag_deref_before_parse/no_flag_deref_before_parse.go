package no_flag_deref_before_parse

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoFlagDerefBeforeParseRule detects immediate dereferencing of flag package
// return values. Functions like flag.String, flag.Int, etc. return pointers
// that should only be dereferenced after flag.Parse() is called.
type NoFlagDerefBeforeParseRule struct{}

// flagFunctions lists the flag package functions that return pointers.
var flagFunctions = []string{
	"Bool",
	"Duration",
	"Float64",
	"Int",
	"Int64",
	"String",
	"Uint",
	"Uint64",
}

// Apply applies the rule to given file.
func (r *NoFlagDerefBeforeParseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintFlagDeref{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoFlagDerefBeforeParseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintFlagDeref{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoFlagDerefBeforeParseRule) Name() string {
	return "noFlagDerefBeforeParse"
}

// Group returns the rule group.
func (*NoFlagDerefBeforeParseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoFlagDerefBeforeParseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintFlagDeref struct {
	onFailure func(lint.Failure)
}

func (w *lintFlagDeref) Visit(node ast.Node) ast.Visitor {
	starExpr, ok := node.(*ast.StarExpr)
	if !ok {
		return w
	}

	callExpr, ok := starExpr.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fn := range flagFunctions {
		if astutils.IsPkgDotName(callExpr.Fun, "flag", fn) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       node,
				Category:   lint.FailureCategoryLogic,
				Failure:    "immediate dereference of flag." + fn + " result before flag.Parse",
			})
			return w
		}
	}

	return w
}
