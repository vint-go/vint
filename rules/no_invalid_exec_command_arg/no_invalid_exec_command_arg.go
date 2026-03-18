package no_invalid_exec_command_arg

import (
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInvalidExecCommandArgRule flags calls to exec.Command where the first
// argument is a string literal containing a space, which likely means the
// caller is passing a shell command instead of just the executable path.
type NoInvalidExecCommandArgRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidExecCommandArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoInvalidExecCommandArg{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidExecCommandArgRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoInvalidExecCommandArg{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidExecCommandArgRule) Name() string {
	return "noInvalidExecCommandArg"
}

// Group returns the rule group.
func (*NoInvalidExecCommandArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidExecCommandArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoInvalidExecCommandArg struct {
	onFailure func(lint.Failure)
}

func (w *lintNoInvalidExecCommandArg) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "exec", "Command") &&
		!astutils.IsPkgDotName(ce.Fun, "exec", "CommandContext") {
		return w
	}

	// For exec.Command the first arg is the executable path (index 0).
	// For exec.CommandContext the first arg is context, executable is index 1.
	argIdx := 0
	if astutils.IsPkgDotName(ce.Fun, "exec", "CommandContext") {
		argIdx = 1
	}

	if len(ce.Args) <= argIdx {
		return w
	}

	arg := ce.Args[argIdx]
	lit, ok := arg.(*ast.BasicLit)
	if !ok {
		return w
	}

	// BasicLit.Value includes quotes, e.g. `"echo hello"`
	val := lit.Value
	if len(val) < 2 {
		return w
	}
	// Strip quotes
	inner := val[1 : len(val)-1]
	if strings.Contains(inner, " ") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "first argument to exec.Command looks like a shell command, but exec.Command expects the executable path without arguments",
		})
	}

	return w
}
