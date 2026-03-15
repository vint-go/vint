package no_variable_command_execution

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoVariableCommandExecutionRule detects calls to command execution functions
// (os/exec.Command, os/exec.CommandContext, syscall.Exec, syscall.ForkExec,
// syscall.StartProcess) where command name or arguments are derived from
// variables rather than constants.
type NoVariableCommandExecutionRule struct{}

// Apply applies the rule to given file.
func (r *NoVariableCommandExecutionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoVariableCommandExecution{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoVariableCommandExecutionRule) Name() string {
	return "noVariableCommandExecution"
}

// Group returns the rule group.
func (*NoVariableCommandExecutionRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoVariableCommandExecutionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoVariableCommandExecution struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

// commandFuncInfo describes a command execution function to check.
type commandFuncInfo struct {
	pkg      string
	name     string
	// argStart is the index of the first argument that is the command name.
	// For exec.CommandContext, the first arg is a context, so command starts at index 1.
	argStart int
}

// commandFuncs lists the command execution functions to check.
var commandFuncs = []commandFuncInfo{
	{pkg: "exec", name: "Command", argStart: 0},
	{pkg: "exec", name: "CommandContext", argStart: 1},
	{pkg: "syscall", name: "Exec", argStart: 0},
	{pkg: "syscall", name: "ForkExec", argStart: 0},
	{pkg: "syscall", name: "StartProcess", argStart: 0},
}

func (w *lintNoVariableCommandExecution) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fn := range commandFuncs {
		if !astutils.IsPkgDotName(ce.Fun, fn.pkg, fn.name) {
			continue
		}

		// Check all arguments starting from argStart.
		for i := fn.argStart; i < len(ce.Args); i++ {
			if !isConstantExpr(ce.Args[i]) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       ce,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "potential command injection: variable argument passed to " + fn.pkg + "." + fn.name,
				})
				break // report once per call
			}
		}
		return w
	}

	return w
}

// isConstantExpr returns true if the expression is a compile-time constant
// (string literal, integer literal, or a constant composed of such).
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// String literals, integer literals, etc.
		return e.Kind == token.STRING || e.Kind == token.INT || e.Kind == token.FLOAT || e.Kind == token.CHAR
	case *ast.BinaryExpr:
		// Constant expressions like "a" + "b"
		return isConstantExpr(e.X) && isConstantExpr(e.Y)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	case *ast.UnaryExpr:
		return isConstantExpr(e.X)
	}
	return false
}
