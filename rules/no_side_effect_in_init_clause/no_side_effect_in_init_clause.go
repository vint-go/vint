package no_side_effect_in_init_clause

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSideEffectInInitClauseRule detects non-assignment statements inside
// if/switch init clauses, recommending they be moved before the conditional.
type NoSideEffectInInitClauseRule struct{}

// Apply applies the rule to given file.
func (r *NoSideEffectInInitClauseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintSideEffectInit{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSideEffectInInitClauseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintSideEffectInit{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSideEffectInInitClauseRule) Name() string {
	return "noSideEffectInInitClause"
}

// Group returns the rule group.
func (*NoSideEffectInInitClauseRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoSideEffectInInitClauseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSideEffectInit struct {
	onFailure func(lint.Failure)
}

func (w *lintSideEffectInit) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.IfStmt:
		w.checkInit(stmt.Init)
	case *ast.SwitchStmt:
		w.checkInit(stmt.Init)
	}
	return w
}

// checkInit reports a failure if the init statement is not an assignment.
func (w *lintSideEffectInit) checkInit(init ast.Stmt) {
	if init == nil {
		return
	}

	// Assignment statements are fine in init clauses
	if _, ok := init.(*ast.AssignStmt); ok {
		return
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       init,
		Category:   lint.FailureCategoryStyle,
		Failure:    "avoid side effects in init clause of conditional statement; move the statement before the if or switch",
	})
}
