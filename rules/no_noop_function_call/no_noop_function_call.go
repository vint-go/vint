package no_noop_function_call

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNoopFunctionCallRule detects function calls with arguments that make the call a no-op.
type NoNoopFunctionCallRule struct{}

// noopFuncCheck describes a function whose last argument of 0 makes it a no-op.
type noopFuncCheck struct {
	pkg      string
	name     string
	argIndex int    // index of the argument to check for 0
	reason   string // human-readable explanation
}

// noopFuncs is the list of standard library functions where passing 0 for the
// count/n argument makes the call effectively a no-op.
var noopFuncs = []noopFuncCheck{
	{pkg: "strings", name: "Replace", argIndex: 3, reason: "replaces 0 occurrences (no-op)"},
	{pkg: "strings", name: "SplitN", argIndex: 2, reason: "returns nil"},
	{pkg: "strings", name: "SplitAfterN", argIndex: 2, reason: "returns nil"},
	{pkg: "bytes", name: "Replace", argIndex: 3, reason: "replaces 0 occurrences (no-op)"},
	{pkg: "bytes", name: "SplitN", argIndex: 2, reason: "returns nil"},
	{pkg: "bytes", name: "SplitAfterN", argIndex: 2, reason: "returns nil"},
}

// Apply applies the rule to given file.
func (r *NoNoopFunctionCallRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoopFunctionCall{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNoopFunctionCallRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoopFunctionCall{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNoopFunctionCallRule) Name() string {
	return "noNoopFunctionCall"
}

// Group returns the rule group.
func (*NoNoopFunctionCallRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNoopFunctionCallRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoopFunctionCall struct {
	onFailure func(lint.Failure)
}

func (w *lintNoopFunctionCall) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, check := range noopFuncs {
		if !astutils.IsPkgDotName(ce.Fun, check.pkg, check.name) {
			continue
		}

		if len(ce.Args) <= check.argIndex {
			continue
		}

		arg := ce.Args[check.argIndex]
		if isIntLiteralZero(arg) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("%s.%s called with n == 0: %s", check.pkg, check.name, check.reason),
			})
		}

		return w
	}

	return w
}

// isIntLiteralZero checks whether an expression is the integer literal 0.
func isIntLiteralZero(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.INT && lit.Value == "0"
}
