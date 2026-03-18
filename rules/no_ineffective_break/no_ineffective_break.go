package no_ineffective_break

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIneffectiveBreakRule detects break statements inside switch/select that
// have no effect because they only break out of the switch/select, not an
// enclosing for loop. If the intent was to break the loop, a labeled break
// should be used instead.
type NoIneffectiveBreakRule struct{}

// Apply applies the rule to given file.
func (r *NoIneffectiveBreakRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoIneffectiveBreak{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIneffectiveBreakRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoIneffectiveBreak{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoIneffectiveBreakRule) Name() string {
	return "noIneffectiveBreak"
}

// Group returns the rule group.
func (*NoIneffectiveBreakRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoIneffectiveBreakRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoIneffectiveBreak struct {
	onFailure func(lint.Failure)
}

func (w *lintNoIneffectiveBreak) Visit(node ast.Node) ast.Visitor {
	// Look for for/range loops that contain switch/select with unlabeled breaks
	switch v := node.(type) {
	case *ast.ForStmt:
		w.checkLoopBody(v.Body)
		return nil // we manually walk the body
	case *ast.RangeStmt:
		w.checkLoopBody(v.Body)
		return nil // we manually walk the body
	}
	return w
}

// checkLoopBody walks the body of a for/range loop looking for switch/select
// statements that contain unlabeled break statements.
func (w *lintNoIneffectiveBreak) checkLoopBody(body *ast.BlockStmt) {
	if body == nil {
		return
	}
	walker := &loopBodyWalker{onFailure: w.onFailure}
	ast.Walk(walker, body)
}

// loopBodyWalker walks inside a loop body looking for switch/select statements
// that contain unlabeled break statements. It stops descending into nested
// loops since a break in a nested loop's switch/select would affect that inner
// loop, not the outer one.
type loopBodyWalker struct {
	onFailure func(lint.Failure)
}

func (w *loopBodyWalker) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.ForStmt:
		// Nested loop: breaks inside here would be about this inner loop.
		// We still need to check for switch/select inside this nested loop.
		innerWalker := &lintNoIneffectiveBreak{onFailure: w.onFailure}
		ast.Walk(innerWalker, v)
		return nil
	case *ast.RangeStmt:
		// Same as above for range loops.
		innerWalker := &lintNoIneffectiveBreak{onFailure: w.onFailure}
		ast.Walk(innerWalker, v)
		return nil
	case *ast.SwitchStmt:
		w.checkSwitchForIneffectiveBreaks(v)
		return nil
	case *ast.TypeSwitchStmt:
		w.checkTypeSwitchForIneffectiveBreaks(v)
		return nil
	case *ast.SelectStmt:
		w.checkSelectForIneffectiveBreaks(v)
		return nil
	}
	return w
}

func (w *loopBodyWalker) checkSwitchForIneffectiveBreaks(sw *ast.SwitchStmt) {
	if sw.Body == nil {
		return
	}
	for _, item := range sw.Body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		w.checkStmtsForUnlabeledBreak(cc.Body)
	}
}

func (w *loopBodyWalker) checkTypeSwitchForIneffectiveBreaks(sw *ast.TypeSwitchStmt) {
	if sw.Body == nil {
		return
	}
	for _, item := range sw.Body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		w.checkStmtsForUnlabeledBreak(cc.Body)
	}
}

func (w *loopBodyWalker) checkSelectForIneffectiveBreaks(sel *ast.SelectStmt) {
	if sel.Body == nil {
		return
	}
	for _, item := range sel.Body.List {
		cc, ok := item.(*ast.CommClause)
		if !ok {
			continue
		}
		w.checkStmtsForUnlabeledBreak(cc.Body)
	}
}

// checkStmtsForUnlabeledBreak walks through the statements in a case clause
// body and reports any unlabeled break statements that only break out of the
// switch/select and not the enclosing loop.
func (w *loopBodyWalker) checkStmtsForUnlabeledBreak(stmts []ast.Stmt) {
	finder := &breakFinder{onFailure: w.onFailure}
	for _, stmt := range stmts {
		ast.Walk(finder, stmt)
	}
}

// breakFinder walks AST nodes looking for unlabeled break statements.
// It stops at nested switch/select/for/range boundaries since those would
// capture the break.
type breakFinder struct {
	onFailure func(lint.Failure)
}

func (f *breakFinder) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.BranchStmt:
		if v.Tok == token.BREAK && v.Label == nil {
			f.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       v,
				Failure:    "ineffective break statement. Did you mean to break out of the outer loop?",
			})
		}
		return nil
	case *ast.ForStmt, *ast.RangeStmt:
		// Nested loop: break inside would break this loop, not our switch/select
		return nil
	case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		// Nested switch/select: break inside would break this, not outer
		return nil
	}
	return f
}
