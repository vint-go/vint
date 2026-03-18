package no_empty_branch

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoEmptyBranchRule detects empty bodies in if or else branches.
// An empty body usually indicates incomplete code.
type NoEmptyBranchRule struct{}

// Apply applies the rule to given file.
func (r *NoEmptyBranchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyBranch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEmptyBranchRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyBranch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoEmptyBranchRule) Name() string {
	return "noEmptyBranch"
}

// Group returns the rule group.
func (*NoEmptyBranchRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoEmptyBranchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoEmptyBranch struct {
	onFailure func(lint.Failure)
}

func (w *lintNoEmptyBranch) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Check the if body
	if ifStmt.Body != nil && len(ifStmt.Body.List) == 0 {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       ifStmt.Body,
			Failure:    "empty branch",
		})
	}

	// Check the else body if it is a block (not another if)
	if elseBlock, ok := ifStmt.Else.(*ast.BlockStmt); ok {
		if len(elseBlock.List) == 0 {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       elseBlock,
				Failure:    "empty branch",
			})
		}
	}

	return w
}
