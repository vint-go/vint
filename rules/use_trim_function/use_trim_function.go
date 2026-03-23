package use_trim_function

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTrimFunctionRule detects manual trimming patterns that should use
// strings.TrimPrefix or strings.TrimSuffix instead.
type UseTrimFunctionRule struct{}

// Apply applies the rule to given file.
func (r *UseTrimFunctionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTrimFunction{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTrimFunctionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTrimFunction{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseTrimFunctionRule) Name() string {
	return "useTrimFunction"
}

// Group returns the rule group.
func (*UseTrimFunctionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTrimFunctionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTrimFunction struct {
	onFailure func(lint.Failure)
}

func (w *lintTrimFunction) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Must have no init statement and no else branch
	if ifStmt.Init != nil || ifStmt.Else != nil {
		return w
	}

	// Body must have exactly one statement
	if len(ifStmt.Body.List) != 1 {
		return w
	}

	// The condition must be a call to strings.HasPrefix or strings.HasSuffix
	callExpr, ok := ifStmt.Cond.(*ast.CallExpr)
	if !ok {
		return w
	}

	isPrefix := astutils.IsPkgDotName(callExpr.Fun, "strings", "HasPrefix")
	isSuffix := astutils.IsPkgDotName(callExpr.Fun, "strings", "HasSuffix")
	if !isPrefix && !isSuffix {
		return w
	}

	if len(callExpr.Args) != 2 {
		return w
	}

	strArg := callExpr.Args[0] // the string being checked
	fixArg := callExpr.Args[1] // the prefix/suffix

	// The body statement must be an assignment
	assignStmt, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
	if !ok {
		return w
	}

	if assignStmt.Tok != token.ASSIGN {
		return w
	}

	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return w
	}

	// LHS must be the same variable as the first arg to HasPrefix/HasSuffix
	if astutils.GoFmt(assignStmt.Lhs[0]) != astutils.GoFmt(strArg) {
		return w
	}

	// RHS must be a slice expression on the same variable
	sliceExpr, ok := assignStmt.Rhs[0].(*ast.SliceExpr)
	if !ok {
		return w
	}

	// The slice target must be the same variable
	if astutils.GoFmt(sliceExpr.X) != astutils.GoFmt(strArg) {
		return w
	}

	if isPrefix {
		w.checkPrefix(ifStmt, sliceExpr, fixArg)
	} else {
		w.checkSuffix(ifStmt, sliceExpr, strArg, fixArg)
	}

	return w
}

// checkPrefix checks for: if strings.HasPrefix(s, prefix) { s = s[len(prefix):] }
// The slice must be s[len(prefix):] — i.e., Low = len(prefix), High = nil
func (w *lintTrimFunction) checkPrefix(ifStmt *ast.IfStmt, sliceExpr *ast.SliceExpr, fixArg ast.Expr) {
	if sliceExpr.High != nil {
		return
	}

	low := sliceExpr.Low
	if low == nil {
		return
	}

	// Low must be len(prefix)
	if !isLenCall(low, fixArg) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       ifStmt,
		Failure:    "replace HasPrefix and manual slicing with strings.TrimPrefix",
	})
}

// checkSuffix checks for: if strings.HasSuffix(s, suffix) { s = s[:len(s)-len(suffix)] }
// The slice must be s[:len(s)-len(suffix)] — i.e., Low = nil, High = len(s)-len(suffix)
func (w *lintTrimFunction) checkSuffix(ifStmt *ast.IfStmt, sliceExpr *ast.SliceExpr, strArg, fixArg ast.Expr) {
	if sliceExpr.Low != nil {
		return
	}

	high := sliceExpr.High
	if high == nil {
		return
	}

	// High must be len(s) - len(suffix)
	binExpr, ok := high.(*ast.BinaryExpr)
	if !ok {
		return
	}

	if binExpr.Op != token.SUB {
		return
	}

	// Left side: len(s)
	if !isLenCall(binExpr.X, strArg) {
		return
	}

	// Right side: len(suffix)
	if !isLenCall(binExpr.Y, fixArg) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       ifStmt,
		Failure:    "replace HasSuffix and manual slicing with strings.TrimSuffix",
	})
}

// isLenCall checks if expr is a call to len(arg) where arg matches the expected expression.
func isLenCall(expr ast.Expr, expectedArg ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "len" {
		return false
	}

	if len(call.Args) != 1 {
		return false
	}

	return astutils.GoFmt(call.Args[0]) == astutils.GoFmt(expectedArg)
}
