package no_unnecessary_defer

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnnecessaryDeferRule detects defer statements that are placed immediately
// before a function's return statement and serve no purpose. The deferred call
// could simply be executed directly without the defer mechanism.
type NoUnnecessaryDeferRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryDeferRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDefer{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryDeferRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDefer{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryDeferRule) Name() string {
	return "noUnnecessaryDefer"
}

// Group returns the rule group.
func (*NoUnnecessaryDeferRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryDeferRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnnecessaryDefer struct {
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryDefer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil {
			w.checkBlockStmt(n.Body)
		}
		return nil // don't recurse; we handle nested blocks ourselves
	case *ast.FuncLit:
		if n.Body != nil {
			w.checkBlockStmt(n.Body)
		}
		return nil // don't recurse into function literals
	}
	return w
}

// checkBlockStmt inspects a block of statements looking for a defer immediately
// followed by a return. It also recurses into nested blocks (if/else, switch, etc.)
// but does NOT recurse into nested function literals (they have their own scope).
func (w *lintUnnecessaryDefer) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil || len(block.List) < 2 {
		return
	}

	for i := 0; i < len(block.List)-1; i++ {
		deferStmt, ok := block.List[i].(*ast.DeferStmt)
		if !ok {
			continue
		}
		// Check if the very next statement is a return
		if _, ok := block.List[i+1].(*ast.ReturnStmt); ok {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       deferStmt,
				Category:   lint.FailureCategoryStyle,
				Failure:    "unnecessary defer: right before return",
			})
		}
	}

	// Recurse into nested block-containing statements, but not function literals
	for _, stmt := range block.List {
		w.recurseIntoStmt(stmt)
	}
}

// recurseIntoStmt descends into statements that contain block bodies
// (if, for, switch, select, etc.) to find nested defer-before-return patterns.
func (w *lintUnnecessaryDefer) recurseIntoStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		w.checkBlockStmt(s.Body)
		if s.Else != nil {
			w.recurseIntoStmt(s.Else)
		}
	case *ast.ForStmt:
		w.checkBlockStmt(s.Body)
	case *ast.RangeStmt:
		w.checkBlockStmt(s.Body)
	case *ast.SwitchStmt:
		if s.Body != nil {
			for _, item := range s.Body.List {
				if cc, ok := item.(*ast.CaseClause); ok {
					w.checkStmtList(cc.Body)
				}
			}
		}
	case *ast.TypeSwitchStmt:
		if s.Body != nil {
			for _, item := range s.Body.List {
				if cc, ok := item.(*ast.CaseClause); ok {
					w.checkStmtList(cc.Body)
				}
			}
		}
	case *ast.SelectStmt:
		if s.Body != nil {
			for _, item := range s.Body.List {
				if cc, ok := item.(*ast.CommClause); ok {
					w.checkStmtList(cc.Body)
				}
			}
		}
	case *ast.BlockStmt:
		w.checkBlockStmt(s)
	case *ast.LabeledStmt:
		w.recurseIntoStmt(s.Stmt)
	}
}

// checkStmtList checks a flat list of statements (e.g., case clause bodies)
// that don't have their own *ast.BlockStmt wrapper.
func (w *lintUnnecessaryDefer) checkStmtList(stmts []ast.Stmt) {
	if len(stmts) < 2 {
		return
	}

	for i := 0; i < len(stmts)-1; i++ {
		deferStmt, ok := stmts[i].(*ast.DeferStmt)
		if !ok {
			continue
		}
		if _, ok := stmts[i+1].(*ast.ReturnStmt); ok {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       deferStmt,
				Category:   lint.FailureCategoryStyle,
				Failure:    "unnecessary defer: right before return",
			})
		}
	}

	// Recurse into nested blocks within the statement list
	for _, stmt := range stmts {
		w.recurseIntoStmt(stmt)
	}
}
