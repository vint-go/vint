package no_unnecessary_block

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryBlockRule detects unnecessary braced statement blocks that
// can be removed without changing program semantics.
type NoUnnecessaryBlockRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryBlockRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnnecessaryBlock{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryBlockRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnnecessaryBlock{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryBlockRule) Name() string {
	return "noUnnecessaryBlock"
}

// Group returns the rule group.
func (*NoUnnecessaryBlockRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryBlockRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoUnnecessaryBlock struct {
	onFailure func(lint.Failure)
}

func (w *lintNoUnnecessaryBlock) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BlockStmt:
		// Check each statement in the block's list for unnecessary nested blocks
		w.checkStatementList(n.List)
	case *ast.CaseClause:
		// Check for case clause with a single block statement
		w.checkCaseClause(n)
	case *ast.CommClause:
		// Check for comm clause (select case) with a single block statement
		w.checkCommClause(n)
	}
	return w
}

// checkStatementList looks for *ast.BlockStmt within a statement list
// that contain no variable definitions or declarations.
func (w *lintNoUnnecessaryBlock) checkStatementList(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		block, ok := stmt.(*ast.BlockStmt)
		if !ok {
			continue
		}
		if !containsVarDeclOrShortAssign(block) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       block,
				Failure:    "unnecessary block statement can be removed",
			})
		}
	}
}

// checkCaseClause checks if a case clause contains exactly one statement
// that is a block statement, which is redundant.
func (w *lintNoUnnecessaryBlock) checkCaseClause(cc *ast.CaseClause) {
	if len(cc.Body) != 1 {
		return
	}
	block, ok := cc.Body[0].(*ast.BlockStmt)
	if !ok {
		return
	}
	_ = block
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       block,
		Failure:    "unnecessary block in case clause",
	})
}

// checkCommClause checks if a select comm clause contains exactly one statement
// that is a block statement, which is redundant.
func (w *lintNoUnnecessaryBlock) checkCommClause(cc *ast.CommClause) {
	if len(cc.Body) != 1 {
		return
	}
	block, ok := cc.Body[0].(*ast.BlockStmt)
	if !ok {
		return
	}
	_ = block
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       block,
		Failure:    "unnecessary block in case clause",
	})
}

// containsVarDeclOrShortAssign returns true if the block contains any
// variable declarations (var ...) or short variable declarations (:=).
func containsVarDeclOrShortAssign(block *ast.BlockStmt) bool {
	for _, stmt := range block.List {
		switch s := stmt.(type) {
		case *ast.DeclStmt:
			// var or const or type declaration
			if genDecl, ok := s.Decl.(*ast.GenDecl); ok {
				if genDecl.Tok == token.VAR {
					return true
				}
			}
		case *ast.AssignStmt:
			if s.Tok == token.DEFINE {
				return true
			}
		}
	}
	return false
}
