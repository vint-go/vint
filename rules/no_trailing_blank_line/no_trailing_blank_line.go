package no_trailing_blank_line

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTrailingBlankLineRule detects unnecessary blank lines at the end of
// function bodies, if/else blocks, for loops, switch statements, or case clauses.
type NoTrailingBlankLineRule struct{}

// Apply applies the rule to given file.
func (r *NoTrailingBlankLineRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintTrailingBlankLine{
		onFailure: onFailure,
		file:      file,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTrailingBlankLineRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTrailingBlankLine{onFailure: onFailure, file: file}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTrailingBlankLineRule) Name() string {
	return "noTrailingBlankLine"
}

// Group returns the rule group.
func (*NoTrailingBlankLineRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoTrailingBlankLineRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTrailingBlankLine struct {
	onFailure func(lint.Failure)
	file      *lint.File
}

func (w *lintTrailingBlankLine) Visit(node ast.Node) ast.Visitor {
	if n, ok := node.(*ast.BlockStmt); ok {
		w.checkBlockStmt(n)
	}
	return w
}

// checkBlockStmt checks if there is a blank line before the closing brace of a block.
func (w *lintTrailingBlankLine) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil || len(block.List) == 0 {
		return
	}

	lastStmt := block.List[len(block.List)-1]
	lastStmtEndLine := w.file.ToPosition(lastStmt.End()).Line
	closeBraceLine := w.file.ToPosition(block.Rbrace).Line

	if closeBraceLine > lastStmtEndLine+1 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       block,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary trailing blank line",
		})
	}
}
