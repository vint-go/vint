package no_exec_command

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExecCommandRule disallows calling exec.Command without a context.
// Use exec.CommandContext instead, which accepts a context.Context as its first argument.
type NoExecCommandRule struct{}

// Apply applies the rule to given file.
func (r *NoExecCommandRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoExecCommand{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExecCommandRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoExecCommand{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoExecCommandRule) Name() string {
	return "noExecCommand"
}

// Group returns the rule group.
func (*NoExecCommandRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoExecCommandRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoExecCommand struct {
	onFailure func(lint.Failure)
}

func (w *lintNoExecCommand) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "exec", "Command") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "exec.Command does not accept a context; use exec.CommandContext instead",
	})

	return w
}
