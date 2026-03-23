package no_redundant_nil_loop_check

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantNilLoopCheckRule detects redundant nil checks around range loops.
// Checking `if s != nil` before `for range s` is redundant because ranging
// over a nil slice or map simply does nothing (zero iterations).
type NoRedundantNilLoopCheckRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantNilLoopCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilLoopCheck{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantNilLoopCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantNilLoopCheck{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantNilLoopCheckRule) Name() string {
	return "noRedundantNilLoopCheck"
}

// Group returns the rule group.
func (*NoRedundantNilLoopCheckRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantNilLoopCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantNilLoopCheck struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantNilLoopCheck) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// The if must have no init statement and no else branch
	if ifStmt.Init != nil || ifStmt.Else != nil {
		return w
	}

	// The condition must be `x != nil`
	nilCheckIdent := extractNilCheckIdent(ifStmt.Cond)
	if nilCheckIdent == "" {
		return w
	}

	// The body must contain exactly one statement: a range loop
	if len(ifStmt.Body.List) != 1 {
		return w
	}

	rangeStmt, ok := ifStmt.Body.List[0].(*ast.RangeStmt)
	if !ok {
		return w
	}

	// The range expression must be the same variable as the nil check
	rangeXStr := astutils.GoFmt(rangeStmt.X)
	if rangeXStr != nilCheckIdent {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ifStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("redundant nil check on %s before range loop; ranging over nil is safe", nilCheckIdent),
	})

	return w
}

// extractNilCheckIdent extracts the variable name from an expression of the form `x != nil`.
func extractNilCheckIdent(expr ast.Expr) string {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return ""
	}

	if binExpr.Op != token.NEQ {
		return ""
	}

	nilIdent, ok := binExpr.Y.(*ast.Ident)
	if !ok {
		return ""
	}

	if nilIdent.Name != "nil" {
		return ""
	}

	return astutils.GoFmt(binExpr.X)
}
