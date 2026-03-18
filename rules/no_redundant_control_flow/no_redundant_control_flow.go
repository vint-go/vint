package no_redundant_control_flow

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantControlFlowRule detects redundant control flow statements:
// a return at the end of a function that returns nothing, or a break
// at the end of a case clause.
type NoRedundantControlFlowRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantControlFlowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantControlFlow{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantControlFlowRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintRedundantControlFlow{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantControlFlowRule) Name() string {
	return "noRedundantControlFlow"
}

// Group returns the rule group.
func (*NoRedundantControlFlowRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantControlFlowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantControlFlow struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantControlFlow) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil && n.Type != nil {
			w.checkRedundantReturn(n.Body, n.Type)
		}
	case *ast.FuncLit:
		if n.Body != nil && n.Type != nil {
			w.checkRedundantReturn(n.Body, n.Type)
		}
	case *ast.SwitchStmt:
		if n.Body != nil {
			w.checkRedundantBreakInSwitch(n.Body)
		}
	case *ast.TypeSwitchStmt:
		if n.Body != nil {
			w.checkRedundantBreakInSwitch(n.Body)
		}
	case *ast.SelectStmt:
		if n.Body != nil {
			w.checkRedundantBreakInSelect(n.Body)
		}
	}
	return w
}

// checkRedundantReturn checks if the last statement in a function body is a
// redundant return (i.e. the function has no return values).
func (w *lintRedundantControlFlow) checkRedundantReturn(body *ast.BlockStmt, funcType *ast.FuncType) {
	// Only flag functions with no return values
	if funcType.Results != nil && len(funcType.Results.List) > 0 {
		return
	}

	if len(body.List) == 0 {
		return
	}

	lastStmt := body.List[len(body.List)-1]
	retStmt, ok := lastStmt.(*ast.ReturnStmt)
	if !ok {
		return
	}

	// Only flag bare returns (no return values)
	if len(retStmt.Results) > 0 {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       retStmt,
		Failure:    "redundant return statement",
	})
}

// checkRedundantBreakInSwitch checks for redundant break statements at the end
// of case clauses in switch statements.
func (w *lintRedundantControlFlow) checkRedundantBreakInSwitch(body *ast.BlockStmt) {
	for _, item := range body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}
		w.checkRedundantBreakInCaseBody(cc.Body)
	}
}

// checkRedundantBreakInSelect checks for redundant break statements at the end
// of comm clauses in select statements.
func (w *lintRedundantControlFlow) checkRedundantBreakInSelect(body *ast.BlockStmt) {
	for _, item := range body.List {
		cc, ok := item.(*ast.CommClause)
		if !ok {
			continue
		}
		w.checkRedundantBreakInCaseBody(cc.Body)
	}
}

// checkRedundantBreakInCaseBody checks if the last statement in a case/comm
// clause body is a redundant break (a break with no label).
func (w *lintRedundantControlFlow) checkRedundantBreakInCaseBody(stmts []ast.Stmt) {
	if len(stmts) == 0 {
		return
	}

	lastStmt := stmts[len(stmts)-1]
	branchStmt, ok := lastStmt.(*ast.BranchStmt)
	if !ok {
		return
	}

	if branchStmt.Tok != token.BREAK {
		return
	}

	// Only flag break without label — labeled breaks are not redundant
	if branchStmt.Label != nil {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       branchStmt,
		Failure:    "redundant break statement",
	})
}
