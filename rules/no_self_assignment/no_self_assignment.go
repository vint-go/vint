package no_self_assignment

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSelfAssignmentRule detects useless self-assignments of the form x = x.
type NoSelfAssignmentRule struct{}

// Apply applies the rule to given file.
func (r *NoSelfAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSelfAssignment{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoSelfAssignmentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSelfAssignment{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSelfAssignmentRule) Name() string {
	return "noSelfAssignment"
}

// Group returns the rule group.
func (*NoSelfAssignmentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSelfAssignmentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoSelfAssignment struct {
	onFailure func(lint.Failure)
}

func (w *lintNoSelfAssignment) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Only check plain assignments (=), not := or +=, etc.
	if assign.Tok.String() != "=" {
		return w
	}

	// Check each pair of LHS and RHS expressions
	if len(assign.Lhs) != len(assign.Rhs) {
		return w
	}

	for i := range assign.Lhs {
		lhs := astutils.GoFmt(assign.Lhs[i])
		rhs := astutils.GoFmt(assign.Rhs[i])
		if lhs != "" && lhs == rhs {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       assign,
				Failure:    "useless self-assignment of " + lhs,
			})
		}
	}

	return w
}
