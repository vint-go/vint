package no_select_break_confusion

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSelectBreakConfusionRule detects break statements inside a select inside a
// for loop that only break out of the select, not the loop. A labeled break
// should be used instead.
type NoSelectBreakConfusionRule struct{}

// Apply applies the rule to given file.
func (r *NoSelectBreakConfusionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSelectBreak{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSelectBreakConfusionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintSelectBreak{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSelectBreakConfusionRule) Name() string {
	return "noSelectBreakConfusion"
}

// Group returns the rule group.
func (*NoSelectBreakConfusionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSelectBreakConfusionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSelectBreak struct {
	onFailure func(lint.Failure)
}

func (w *lintSelectBreak) Visit(node ast.Node) ast.Visitor {
	var body *ast.BlockStmt

	switch n := node.(type) {
	case *ast.ForStmt:
		body = n.Body
	case *ast.RangeStmt:
		body = n.Body
	default:
		return w
	}

	if body == nil {
		return w
	}

	w.checkForBody(body)

	return w
}

// checkForBody inspects the body of a for loop for select statements
// that contain unlabeled break statements.
func (w *lintSelectBreak) checkForBody(body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		// Don't descend into nested for/range loops or function literals,
		// because a break in those scopes wouldn't be confusing with
		// respect to the outer for loop.
		switch n.(type) {
		case *ast.ForStmt:
			// This is the outer for we are already inside (on first call)
			// or a nested for. We handle the outer case by starting Inspect
			// on the body, not the for itself. For nested fors, stop.
			return false
		case *ast.RangeStmt:
			return false
		case *ast.FuncLit:
			return false
		}

		selectStmt, ok := n.(*ast.SelectStmt)
		if !ok {
			return true
		}

		// Found a select inside a for. Check its cases for unlabeled breaks.
		w.checkSelectCases(selectStmt)

		// Don't descend further into the select; we've already checked it.
		return false
	})
}

// checkSelectCases looks at every case clause in a select statement and
// reports unlabeled break statements.
func (w *lintSelectBreak) checkSelectCases(sel *ast.SelectStmt) {
	if sel.Body == nil {
		return
	}

	for _, stmt := range sel.Body.List {
		cc, ok := stmt.(*ast.CommClause)
		if !ok {
			continue
		}

		w.findUnlabeledBreaks(cc)
	}
}

// findUnlabeledBreaks walks the statements in a comm clause to find
// unlabeled break statements. It stops descending into nested
// for/range/switch/select/func literals since breaks in those scopes
// target the inner construct.
func (w *lintSelectBreak) findUnlabeledBreaks(cc *ast.CommClause) {
	for _, stmt := range cc.Body {
		w.inspectForBreaks(stmt)
	}
}

func (w *lintSelectBreak) inspectForBreaks(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.FuncLit:
			// Break inside these would target them, not the select.
			return false
		}

		branchStmt, ok := n.(*ast.BranchStmt)
		if !ok {
			return true
		}

		if branchStmt.Tok == token.BREAK && branchStmt.Label == nil {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       branchStmt,
				Category:   lint.FailureCategoryLogic,
				Failure:    "break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop",
			})
		}

		return true
	})
}
