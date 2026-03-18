package use_error_last_return

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseErrorLastReturnRule checks that a function's error value is its last return value.
type UseErrorLastReturnRule struct{}

// Apply applies the rule to given file.
func (r *UseErrorLastReturnRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintErrorLastReturn{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseErrorLastReturnRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintErrorLastReturn{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseErrorLastReturnRule) Name() string {
	return "useErrorLastReturn"
}

// Group returns the rule group.
func (*UseErrorLastReturnRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseErrorLastReturnRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintErrorLastReturn struct {
	onFailure func(lint.Failure)
}

func (w *lintErrorLastReturn) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	results := funcDecl.Type.Results
	if results == nil || len(results.List) < 2 {
		return w
	}

	// Check if the last return value is already error.
	lastResult := results.List[len(results.List)-1]
	if astutils.IsIdent(lastResult.Type, "error") {
		return w
	}

	// Check if any non-last return value is error.
	for _, r := range results.List[:len(results.List)-1] {
		if astutils.IsIdent(r.Type, "error") {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 0.9,
				Node:       funcDecl,
				Failure:    "error should be the last type when returning multiple items",
			})

			break // only flag once per function
		}
	}

	return w
}
