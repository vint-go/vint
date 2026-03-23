package no_duplicate_branch_body

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDuplicateBranchBodyRule detects if statements where both the then-branch
// and else-branch contain identical code bodies.
type NoDuplicateBranchBodyRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateBranchBodyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateBranchBody{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateBranchBodyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateBranchBody{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateBranchBodyRule) Name() string {
	return "noDuplicateBranchBody"
}

// Group returns the rule group.
func (*NoDuplicateBranchBodyRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateBranchBodyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoDuplicateBranchBody struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDuplicateBranchBody) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	if ifStmt.Else == nil {
		return w // if without else
	}

	elseBranch, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok {
		// if-else-if construction; only check simple if...else
		return w
	}

	if len(ifStmt.Body.List) == 0 && len(elseBranch.List) == 0 {
		return w // both branches empty, not interesting
	}

	if len(ifStmt.Body.List) != len(elseBranch.List) {
		return w // different number of statements
	}

	bodyStr := astutils.GoFmt(ifStmt.Body)
	elseStr := astutils.GoFmt(elseBranch)

	if bodyStr == elseStr {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ifStmt,
			Category:   lint.FailureCategoryLogic,
			Failure:    "both branches of the if statement have identical bodies",
		})
	}

	return w
}
