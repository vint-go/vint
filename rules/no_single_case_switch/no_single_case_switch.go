package no_single_case_switch

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSingleCaseSwitchRule detects switch statements that could be better written
// as if statements, including single-case switches and default-only switches.
type NoSingleCaseSwitchRule struct{}

// Apply applies the rule to given file.
func (r *NoSingleCaseSwitchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSingleCaseSwitch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSingleCaseSwitchRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSingleCaseSwitch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSingleCaseSwitchRule) Name() string {
	return "noSingleCaseSwitch"
}

// Group returns the rule group.
func (*NoSingleCaseSwitchRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoSingleCaseSwitchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoSingleCaseSwitch struct {
	onFailure func(lint.Failure)
}

func (w *lintNoSingleCaseSwitch) Visit(node ast.Node) ast.Visitor {
	switchStmt, ok := node.(*ast.SwitchStmt)
	if !ok {
		return w
	}

	if switchStmt.Body == nil {
		return w
	}

	clauses := switchStmt.Body.List
	if len(clauses) != 1 {
		return w
	}

	cc, ok := clauses[0].(*ast.CaseClause)
	if !ok {
		return w
	}

	// Check if the case clause contains a break statement; if so, skip.
	if containsBreak(cc) {
		return w
	}

	if cc.List == nil {
		// Default-only switch
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       switchStmt,
			Failure:    "switch with only a default case is redundant",
		})
	} else {
		// Single case switch (not default)
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       switchStmt,
			Failure:    "switch with a single case can be rewritten as an if statement",
		})
	}

	return w
}

// containsBreak checks whether the case clause body contains a break statement.
func containsBreak(cc *ast.CaseClause) bool {
	for _, stmt := range cc.Body {
		if branchStmt, ok := stmt.(*ast.BranchStmt); ok && branchStmt.Tok == token.BREAK {
			return true
		}
	}
	return false
}
