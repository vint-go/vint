package no_single_iteration_loop

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSingleIterationLoopRule detects loops that exit unconditionally after one iteration.
type NoSingleIterationLoopRule struct{}

// Apply applies the rule to given file.
func (r *NoSingleIterationLoopRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSingleIterLoop{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSingleIterationLoopRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintSingleIterLoop{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSingleIterationLoopRule) Name() string {
	return "noSingleIterationLoop"
}

// Group returns the rule group.
func (*NoSingleIterationLoopRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSingleIterationLoopRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSingleIterLoop struct {
	onFailure func(lint.Failure)
}

func (w *lintSingleIterLoop) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ForStmt:
		// Skip infinite loops (for { ... }) — these are intentional do-once blocks
		if n.Init == nil && n.Cond == nil && n.Post == nil {
			return w
		}
		if n.Body != nil && bodyAlwaysExits(n.Body, "") {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   lint.FailureCategoryLogic,
				Failure:    "loop exits unconditionally after one iteration",
			})
		}
	case *ast.RangeStmt:
		if n.Body != nil && bodyAlwaysExits(n.Body, "") {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   lint.FailureCategoryLogic,
				Failure:    "loop exits unconditionally after one iteration",
			})
		}
	}
	return w
}

// bodyAlwaysExits returns true if the block statement always terminates the
// enclosing loop on the first iteration (via return, break, or continue).
// loopLabel is the label of the current loop (empty string if unlabeled).
func bodyAlwaysExits(body *ast.BlockStmt, loopLabel string) bool {
	return stmtListAlwaysExits(body.List, loopLabel)
}

// stmtListAlwaysExits returns true if the given statement list always terminates
// the enclosing loop before reaching the end (i.e., before a second iteration could happen).
func stmtListAlwaysExits(stmts []ast.Stmt, loopLabel string) bool {
	for _, stmt := range stmts {
		if stmtAlwaysExits(stmt, loopLabel) {
			return true
		}
	}
	return false
}

// stmtAlwaysExits returns true if a single statement always terminates the
// enclosing loop.
func stmtAlwaysExits(stmt ast.Stmt, loopLabel string) bool {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		switch s.Tok {
		case token.BREAK:
			// break without label breaks the current loop
			if s.Label == nil {
				return true
			}
			// break with a label that matches our loop's label breaks our loop
			if loopLabel != "" && s.Label.Name == loopLabel {
				return true
			}
		case token.CONTINUE:
			// continue without label in the loop body means the loop will
			// iterate again, so it does NOT exit the loop.
			// continue with a label matching our loop also re-iterates.
			return false
		}
		return false

	case *ast.IfStmt:
		// If both branches always exit, the whole if always exits
		if s.Else != nil {
			thenExits := stmtListAlwaysExits(s.Body.List, loopLabel)
			elseExits := stmtAlwaysExits(s.Else, loopLabel)
			return thenExits && elseExits
		}
		return false

	case *ast.BlockStmt:
		return stmtListAlwaysExits(s.List, loopLabel)

	case *ast.SwitchStmt:
		return switchAlwaysExits(s.Body, loopLabel)

	case *ast.TypeSwitchStmt:
		return typeSwitchAlwaysExits(s.Body, loopLabel)

	case *ast.SelectStmt:
		return selectAlwaysExits(s.Body, loopLabel)

	case *ast.LabeledStmt:
		// Check if this is a labeled loop — if so, don't flag it as single-iteration
		// (let the inner loop be checked on its own).
		// But do check the inner statement for exits from our enclosing loop.
		return stmtAlwaysExits(s.Stmt, loopLabel)

	case *ast.ExprStmt:
		// Check for panic calls
		if call, ok := s.X.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				return true
			}
		}
		return false
	}
	return false
}

// switchAlwaysExits checks if a switch statement always exits the enclosing loop.
// This requires a default case and all cases to exit.
func switchAlwaysExits(body *ast.BlockStmt, loopLabel string) bool {
	if body == nil {
		return false
	}

	hasDefault := false
	for _, item := range body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		if cc.List == nil {
			hasDefault = true
		}
		if !caseBodyAlwaysExits(cc.Body, loopLabel) {
			return false
		}
	}
	return hasDefault
}

// typeSwitchAlwaysExits checks if a type switch always exits the enclosing loop.
func typeSwitchAlwaysExits(body *ast.BlockStmt, loopLabel string) bool {
	if body == nil {
		return false
	}

	hasDefault := false
	for _, item := range body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		if cc.List == nil {
			hasDefault = true
		}
		if !caseBodyAlwaysExits(cc.Body, loopLabel) {
			return false
		}
	}
	return hasDefault
}

// selectAlwaysExits checks if a select statement always exits the enclosing loop.
func selectAlwaysExits(body *ast.BlockStmt, loopLabel string) bool {
	if body == nil {
		return false
	}

	hasDefault := false
	for _, item := range body.List {
		cc, ok := item.(*ast.CommClause)
		if !ok {
			continue
		}
		if cc.Comm == nil {
			hasDefault = true
		}
		if !stmtListAlwaysExits(cc.Body, loopLabel) {
			return false
		}
	}
	return hasDefault
}

// caseBodyAlwaysExits checks if a case clause body always exits the enclosing loop.
// Note: a break inside a switch breaks the switch, not the enclosing loop,
// so we need special handling.
func caseBodyAlwaysExits(stmts []ast.Stmt, loopLabel string) bool {
	for _, stmt := range stmts {
		if caseStmtAlwaysExitsLoop(stmt, loopLabel) {
			return true
		}
	}
	return false
}

// caseStmtAlwaysExitsLoop checks if a statement inside a switch case
// always exits the enclosing loop. Break statements inside a switch case
// break the switch, not the enclosing loop, unless they have a label
// matching the loop.
func caseStmtAlwaysExitsLoop(stmt ast.Stmt, loopLabel string) bool {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		switch s.Tok {
		case token.BREAK:
			// A break inside a switch case breaks the switch, not the loop.
			// Only a labeled break targeting the loop label exits the loop.
			if s.Label != nil && loopLabel != "" && s.Label.Name == loopLabel {
				return true
			}
			return false
		}
		return false

	case *ast.IfStmt:
		if s.Else != nil {
			thenExits := caseBodyAlwaysExits(s.Body.List, loopLabel)
			elseExits := caseStmtAlwaysExitsLoop(s.Else, loopLabel)
			return thenExits && elseExits
		}
		return false

	case *ast.BlockStmt:
		return caseBodyAlwaysExits(s.List, loopLabel)

	case *ast.ExprStmt:
		if call, ok := s.X.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				return true
			}
		}
		return false
	}
	return false
}
