package no_leading_blank_line

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoLeadingBlankLineRule detects unnecessary blank lines at the beginning of
// function bodies, if/else blocks, for loops, switch statements, or case clauses.
type NoLeadingBlankLineRule struct{}

// Apply applies the rule to given file.
func (r *NoLeadingBlankLineRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintLeadingBlankLine{
		onFailure: onFailure,
		file:      file,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoLeadingBlankLineRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintLeadingBlankLine{onFailure: onFailure, file: file}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoLeadingBlankLineRule) Name() string {
	return "noLeadingBlankLine"
}

// Group returns the rule group.
func (*NoLeadingBlankLineRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoLeadingBlankLineRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintLeadingBlankLine struct {
	onFailure func(lint.Failure)
	file      *lint.File
}

func (w *lintLeadingBlankLine) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BlockStmt:
		w.checkBlockStmt(n)
	case *ast.CaseClause:
		w.checkCaseClause(n)
	case *ast.CommClause:
		w.checkCommClause(n)
	}
	return w
}

// checkBlockStmt checks if there is a blank line after the opening brace of a block.
func (w *lintLeadingBlankLine) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil || len(block.List) == 0 {
		return
	}

	openBraceLine := w.file.ToPosition(block.Lbrace).Line
	firstStmtLine := w.file.ToPosition(block.List[0].Pos()).Line

	if firstStmtLine > openBraceLine+1 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       block,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary leading blank line",
		})
	}
}

// checkCaseClause checks if there is a blank line after the colon of a case clause.
func (w *lintLeadingBlankLine) checkCaseClause(cc *ast.CaseClause) {
	if len(cc.Body) == 0 {
		return
	}

	colonLine := w.file.ToPosition(cc.Colon).Line
	firstStmtLine := w.file.ToPosition(cc.Body[0].Pos()).Line

	if firstStmtLine > colonLine+1 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cc,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary leading blank line",
		})
	}
}

// checkCommClause checks if there is a blank line after the colon of a comm clause (select case).
func (w *lintLeadingBlankLine) checkCommClause(cc *ast.CommClause) {
	if len(cc.Body) == 0 {
		return
	}

	colonLine := w.file.ToPosition(cc.Colon).Line
	firstStmtLine := w.file.ToPosition(cc.Body[0].Pos()).Line

	if firstStmtLine > colonLine+1 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       cc,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary leading blank line",
		})
	}
}
