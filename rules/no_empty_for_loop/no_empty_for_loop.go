package no_empty_for_loop

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoEmptyForLoopRule detects empty for loops that busy-wait and consume 100% CPU.
type NoEmptyForLoopRule struct{}

// Apply applies the rule to given file.
func (r *NoEmptyForLoopRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyForLoop{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEmptyForLoopRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyForLoop{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoEmptyForLoopRule) Name() string {
	return "noEmptyForLoop"
}

// Group returns the rule group.
func (*NoEmptyForLoopRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoEmptyForLoopRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoEmptyForLoop struct {
	onFailure func(lint.Failure)
}

func (w *lintNoEmptyForLoop) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Only flag infinite for loops: no init, no condition, no post
	if forStmt.Init != nil || forStmt.Cond != nil || forStmt.Post != nil {
		return w
	}

	// Check if the body is empty (no statements)
	if forStmt.Body == nil || len(forStmt.Body.List) == 0 {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       forStmt,
			Failure:    "empty for loop busy-waits and consumes 100% CPU; use select, channel, or runtime.Gosched()",
		})
	}

	return w
}
