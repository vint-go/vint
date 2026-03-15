package no_duplicate_argument

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateArgumentRule detects suspicious duplicated arguments in function calls.
type NoDuplicateArgumentRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateArgumentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateArgument{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateArgumentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateArgument{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateArgumentRule) Name() string {
	return "noDuplicateArgument"
}

// Group returns the rule group.
func (*NoDuplicateArgumentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateArgumentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoDuplicateArgument struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDuplicateArgument) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if len(ce.Args) < 2 {
		return w
	}

	// Check each pair of arguments for duplicates
	seen := make(map[string]int) // arg text -> first index
	for i, arg := range ce.Args {
		argStr := astutils.GoFmt(arg)
		if argStr == "" {
			continue
		}
		if firstIdx, exists := seen[argStr]; exists {
			_ = firstIdx // we know there is a duplicate
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       ce,
				Failure:    "suspicious duplicated argument: " + argStr + " appears more than once in the same call",
			})
			break // report once per call
		}
		seen[argStr] = i
	}

	return w
}
