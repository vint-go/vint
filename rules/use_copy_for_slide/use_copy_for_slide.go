package use_copy_for_slide

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseCopyForSlideRule detects for loops that shift elements within a slice
// and suggests using the built-in copy function instead.
type UseCopyForSlideRule struct{}

// Apply applies the rule to given file.
func (r *UseCopyForSlideRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCopyForSlide{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseCopyForSlideRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCopyForSlide{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseCopyForSlideRule) Name() string {
	return "useCopyForSlide"
}

// Group returns the rule group.
func (*UseCopyForSlideRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseCopyForSlideRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCopyForSlide struct {
	onFailure func(lint.Failure)
}

func (w *lintCopyForSlide) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Must have init, cond, post, and a body with exactly one statement.
	if forStmt.Init == nil || forStmt.Cond == nil || forStmt.Post == nil {
		return w
	}
	if len(forStmt.Body.List) != 1 {
		return w
	}

	// Init: i := 0
	initAssign, ok := forStmt.Init.(*ast.AssignStmt)
	if !ok || initAssign.Tok != token.DEFINE {
		return w
	}
	if len(initAssign.Lhs) != 1 || len(initAssign.Rhs) != 1 {
		return w
	}
	loopVar, ok := initAssign.Lhs[0].(*ast.Ident)
	if !ok {
		return w
	}
	// RHS must be 0
	initVal, ok := initAssign.Rhs[0].(*ast.BasicLit)
	if !ok || initVal.Kind != token.INT || initVal.Value != "0" {
		return w
	}

	// Cond: i < len(s) - 1
	condExpr, ok := forStmt.Cond.(*ast.BinaryExpr)
	if !ok || condExpr.Op != token.LSS {
		return w
	}
	condLhs, ok := condExpr.X.(*ast.Ident)
	if !ok || condLhs.Name != loopVar.Name {
		return w
	}

	// RHS of condition: len(s) - 1
	condRhs, ok := condExpr.Y.(*ast.BinaryExpr)
	if !ok || condRhs.Op != token.SUB {
		return w
	}
	lenCall, ok := condRhs.X.(*ast.CallExpr)
	if !ok {
		return w
	}
	lenIdent, ok := lenCall.Fun.(*ast.Ident)
	if !ok || lenIdent.Name != "len" || len(lenCall.Args) != 1 {
		return w
	}
	sliceStr := astutils.GoFmt(lenCall.Args[0])
	if sliceStr == "" {
		return w
	}
	subVal, ok := condRhs.Y.(*ast.BasicLit)
	if !ok || subVal.Kind != token.INT || subVal.Value != "1" {
		return w
	}

	// Post: i++
	postInc, ok := forStmt.Post.(*ast.IncDecStmt)
	if !ok || postInc.Tok != token.INC {
		return w
	}
	postIdent, ok := postInc.X.(*ast.Ident)
	if !ok || postIdent.Name != loopVar.Name {
		return w
	}

	// Body: s[i] = s[i+1]
	bodyAssign, ok := forStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || bodyAssign.Tok != token.ASSIGN {
		return w
	}
	if len(bodyAssign.Lhs) != 1 || len(bodyAssign.Rhs) != 1 {
		return w
	}

	// LHS: s[i]
	lhsIndex, ok := bodyAssign.Lhs[0].(*ast.IndexExpr)
	if !ok {
		return w
	}
	lhsSliceStr := astutils.GoFmt(lhsIndex.X)
	if lhsSliceStr != sliceStr {
		return w
	}
	lhsIdxIdent, ok := lhsIndex.Index.(*ast.Ident)
	if !ok || lhsIdxIdent.Name != loopVar.Name {
		return w
	}

	// RHS: s[i+1]
	rhsIndex, ok := bodyAssign.Rhs[0].(*ast.IndexExpr)
	if !ok {
		return w
	}
	rhsSliceStr := astutils.GoFmt(rhsIndex.X)
	if rhsSliceStr != sliceStr {
		return w
	}
	rhsAdd, ok := rhsIndex.Index.(*ast.BinaryExpr)
	if !ok || rhsAdd.Op != token.ADD {
		return w
	}
	rhsIdxIdent, ok := rhsAdd.X.(*ast.Ident)
	if !ok || rhsIdxIdent.Name != loopVar.Name {
		return w
	}
	rhsAddVal, ok := rhsAdd.Y.(*ast.BasicLit)
	if !ok || rhsAddVal.Kind != token.INT || rhsAddVal.Value != "1" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       forStmt,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("should replace loop with copy(%s, %s[1:])", sliceStr, sliceStr),
	})

	return w
}
