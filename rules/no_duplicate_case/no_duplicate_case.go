package no_duplicate_case

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateCaseRule detects duplicated case clauses inside switch or select statements.
type NoDuplicateCaseRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateCaseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateCase{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateCaseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDuplicateCase{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateCaseRule) Name() string {
	return "noDuplicateCase"
}

// Group returns the rule group.
func (*NoDuplicateCaseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateCaseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoDuplicateCase struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDuplicateCase) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.SwitchStmt:
		w.checkSwitchStmt(stmt)
	case *ast.TypeSwitchStmt:
		w.checkTypeSwitchStmt(stmt)
	case *ast.SelectStmt:
		w.checkSelectStmt(stmt)
	}
	return w
}

func (w *lintNoDuplicateCase) checkSwitchStmt(stmt *ast.SwitchStmt) {
	seen := map[string]*ast.CaseClause{}
	for _, item := range stmt.Body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		for _, expr := range cc.List {
			key := astutils.GoFmt(expr)
			if key == "" {
				continue
			}
			if _, exists := seen[key]; exists {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       cc,
					Failure:    fmt.Sprintf("duplicate case %s in switch statement", key),
				})
			} else {
				seen[key] = cc
			}
		}
	}
}

func (w *lintNoDuplicateCase) checkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
	seen := map[string]*ast.CaseClause{}
	for _, item := range stmt.Body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		for _, expr := range cc.List {
			key := astutils.GoFmt(expr)
			if key == "" {
				continue
			}
			if _, exists := seen[key]; exists {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       cc,
					Failure:    fmt.Sprintf("duplicate case %s in switch statement", key),
				})
			} else {
				seen[key] = cc
			}
		}
	}
}

func (w *lintNoDuplicateCase) checkSelectStmt(stmt *ast.SelectStmt) {
	seen := map[string]*ast.CommClause{}
	for _, item := range stmt.Body.List {
		cc, ok := item.(*ast.CommClause)
		if !ok {
			continue
		}
		if cc.Comm == nil {
			// default case
			continue
		}
		key := astutils.GoFmt(cc.Comm)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       cc,
				Failure:    fmt.Sprintf("duplicate case %s in select statement", key),
			})
		} else {
			seen[key] = cc
		}
	}
}
