package use_tagged_switch

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseTaggedSwitchRule detects untagged switch statements where all case clauses
// compare against the same variable using ==, and suggests converting them to
// tagged switch statements.
type UseTaggedSwitchRule struct{}

// Apply applies the rule to given file.
func (r *UseTaggedSwitchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseTaggedSwitch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTaggedSwitchRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseTaggedSwitch{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTaggedSwitchRule) Name() string {
	return "useTaggedSwitch"
}

// Group returns the rule group.
func (*UseTaggedSwitchRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTaggedSwitchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseTaggedSwitch struct {
	onFailure func(lint.Failure)
}

func (w *lintUseTaggedSwitch) Visit(node ast.Node) ast.Visitor {
	switchStmt, ok := node.(*ast.SwitchStmt)
	if !ok {
		return w
	}

	// Only check untagged switch statements (switch { ... })
	if switchStmt.Tag != nil {
		return w
	}

	if switchStmt.Body == nil || len(switchStmt.Body.List) == 0 {
		return w
	}

	// Collect all non-default case clauses
	var caseClauses []*ast.CaseClause
	for _, stmt := range switchStmt.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			return w
		}
		if cc.List != nil {
			caseClauses = append(caseClauses, cc)
		}
	}

	// Need at least one non-default case clause to detect a pattern
	if len(caseClauses) == 0 {
		return w
	}

	// Check if all case expressions are == comparisons against the same variable
	var commonVar string
	for _, cc := range caseClauses {
		for _, expr := range cc.List {
			binExpr, ok := expr.(*ast.BinaryExpr)
			if !ok {
				return w
			}
			if binExpr.Op != token.EQL {
				return w
			}

			varStr := astutils.GoFmt(binExpr.X)
			if commonVar == "" {
				commonVar = varStr
			} else if varStr != commonVar {
				return w
			}
		}
	}

	if commonVar == "" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryStyle,
		Failure:    "could convert untagged switch to tagged switch on " + commonVar,
		Node:       switchStmt,
	})

	return w
}
