package use_time_sleep

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseTimeSleepRule detects elaborate sleep patterns using
// select { case <-time.After(d): } and suggests using time.Sleep(d) instead.
type UseTimeSleepRule struct{}

// Apply applies the rule to given file.
func (r *UseTimeSleepRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeSleep{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTimeSleepRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeSleep{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseTimeSleepRule) Name() string {
	return "useTimeSleep"
}

// Group returns the rule group.
func (*UseTimeSleepRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTimeSleepRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTimeSleep struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeSleep) Visit(node ast.Node) ast.Visitor {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return w
	}

	if selectStmt.Body == nil {
		return w
	}

	clauses := selectStmt.Body.List
	if len(clauses) != 1 {
		return w
	}

	cc, ok := clauses[0].(*ast.CommClause)
	if !ok {
		return w
	}

	// A default clause has cc.Comm == nil; skip those.
	if cc.Comm == nil {
		return w
	}

	// The comm must be a receive from time.After(...).
	// It can be either:
	//   case <-time.After(d):       (ExprStmt with UnaryExpr)
	//   case x := <-time.After(d):  (AssignStmt with UnaryExpr on RHS)
	var recvExpr *ast.UnaryExpr

	switch comm := cc.Comm.(type) {
	case *ast.ExprStmt:
		if unary, ok := comm.X.(*ast.UnaryExpr); ok {
			recvExpr = unary
		}
	case *ast.AssignStmt:
		if len(comm.Rhs) == 1 {
			if unary, ok := comm.Rhs[0].(*ast.UnaryExpr); ok {
				recvExpr = unary
			}
		}
	}

	if recvExpr == nil {
		return w
	}

	// Check that it's a channel receive (<-)
	if recvExpr.Op.String() != "<-" {
		return w
	}

	// Check that the operand is a call to time.After
	callExpr, ok := recvExpr.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(callExpr.Fun, "time", "After") {
		return w
	}

	// The case body must be empty (no statements) for a pure sleep equivalent.
	if len(cc.Body) != 0 {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       selectStmt,
		Failure:    "use time.Sleep instead of select{case <-time.After(...)}",
	})

	return w
}
