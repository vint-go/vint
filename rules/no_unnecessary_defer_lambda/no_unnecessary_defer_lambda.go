package no_unnecessary_defer_lambda

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnnecessaryDeferLambdaRule detects deferred function literals that can be
// simplified by deferring the inner function call directly.
type NoUnnecessaryDeferLambdaRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryDeferLambdaRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDeferLambda{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryDeferLambdaRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDeferLambda{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryDeferLambdaRule) Name() string {
	return "noUnnecessaryDeferLambda"
}

// Group returns the rule group.
func (*NoUnnecessaryDeferLambdaRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryDeferLambdaRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnnecessaryDeferLambda struct {
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryDeferLambda) Visit(node ast.Node) ast.Visitor {
	deferStmt, ok := node.(*ast.DeferStmt)
	if !ok {
		return w
	}

	// The deferred call must be a call to a function literal: defer func() { ... }()
	callExpr := deferStmt.Call
	funcLit, ok := callExpr.Fun.(*ast.FuncLit)
	if !ok {
		return w
	}

	// The function literal must take no parameters and return no results
	if funcLit.Type.Params != nil && len(funcLit.Type.Params.List) > 0 {
		return w
	}
	if funcLit.Type.Results != nil && len(funcLit.Type.Results.List) > 0 {
		return w
	}

	// The outer call must pass no arguments
	if len(callExpr.Args) > 0 {
		return w
	}

	// The function literal body must contain exactly one statement
	if funcLit.Body == nil || len(funcLit.Body.List) != 1 {
		return w
	}

	// That single statement must be an expression statement containing a call expression
	exprStmt, ok := funcLit.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return w
	}
	_, ok = exprStmt.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       deferStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    "unnecessary defer lambda, call the function directly",
	})

	return w
}
