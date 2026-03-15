package no_slice_bounds_out_of_range

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/lint"
)

// NoSliceBoundsOutOfRangeRule detects possible slice bounds out of range errors.
type NoSliceBoundsOutOfRangeRule struct{}

// Apply applies the rule to given file.
func (r *NoSliceBoundsOutOfRangeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSliceBoundsOutOfRange{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		w.checkFunctionBody(funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoSliceBoundsOutOfRangeRule) Name() string {
	return "noSliceBoundsOutOfRange"
}

// Group returns the rule group.
func (*NoSliceBoundsOutOfRangeRule) Group() string {
	return "correctness"
}

func (*NoSliceBoundsOutOfRangeRule) RequiresTypecheck() bool {
	return true
}

type lintNoSliceBoundsOutOfRange struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSliceBoundsOutOfRange) checkFunctionBody(body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.IndexExpr:
			w.checkIndexExpr(expr, body)
		case *ast.SliceExpr:
			w.checkSliceExpr(expr, body)
		}
		return true
	})
}

func (w *lintNoSliceBoundsOutOfRange) checkIndexExpr(expr *ast.IndexExpr, body *ast.BlockStmt) {
	if !w.isSliceType(expr.X) {
		return
	}

	sliceName := w.exprName(expr.X)
	if sliceName == "" {
		return
	}

	if w.hasBoundsCheck(body, sliceName, expr.Pos()) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 0.8,
		Node:       expr,
		Failure:    "possible slice bounds out of range",
	})
}

func (w *lintNoSliceBoundsOutOfRange) checkSliceExpr(expr *ast.SliceExpr, body *ast.BlockStmt) {
	if !w.isSliceType(expr.X) {
		return
	}

	sliceName := w.exprName(expr.X)
	if sliceName == "" {
		return
	}

	// Check if high bound is a literal or variable without bounds check
	if expr.High == nil && expr.Low == nil {
		return // s[:] is always safe
	}

	if w.hasBoundsCheck(body, sliceName, expr.Pos()) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 0.8,
		Node:       expr,
		Failure:    "possible slice bounds out of range",
	})
}

func (w *lintNoSliceBoundsOutOfRange) isSliceType(expr ast.Expr) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}
	// Dereference named types
	t = t.Underlying()
	_, ok := t.(*types.Slice)
	return ok
}

func (w *lintNoSliceBoundsOutOfRange) exprName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		prefix := w.exprName(e.X)
		if prefix != "" {
			return prefix + "." + e.Sel.Name
		}
		return e.Sel.Name
	}
	return ""
}

// hasBoundsCheck checks if there's a len() bounds check for the slice
// before the given position within the enclosing function body.
func (w *lintNoSliceBoundsOutOfRange) hasBoundsCheck(body *ast.BlockStmt, sliceName string, pos token.Pos) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}

		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check if this if statement has a len() check for our slice
		if w.conditionChecksLen(ifStmt.Cond, sliceName) {
			// Check if the slice access is inside this if's body or else
			if containsPos(ifStmt, pos) {
				found = true
				return false
			}
		}

		return true
	})
	return found
}

// conditionChecksLen returns true if the condition involves len(sliceName).
func (w *lintNoSliceBoundsOutOfRange) conditionChecksLen(cond ast.Expr, sliceName string) bool {
	switch expr := cond.(type) {
	case *ast.BinaryExpr:
		// Check for comparisons like len(s) > 0, idx < len(s), len(s) >= 5, etc.
		switch expr.Op {
		case token.GTR, token.GEQ, token.LSS, token.LEQ, token.NEQ, token.EQL:
			if w.isLenCall(expr.X, sliceName) || w.isLenCall(expr.Y, sliceName) {
				return true
			}
		case token.LAND, token.LOR:
			return w.conditionChecksLen(expr.X, sliceName) || w.conditionChecksLen(expr.Y, sliceName)
		}
	case *ast.ParenExpr:
		return w.conditionChecksLen(expr.X, sliceName)
	case *ast.UnaryExpr:
		if expr.Op == token.NOT {
			return w.conditionChecksLen(expr.X, sliceName)
		}
	}
	return false
}

// isLenCall checks if the expression is a call to len(sliceName).
func (w *lintNoSliceBoundsOutOfRange) isLenCall(expr ast.Expr, sliceName string) bool {
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

	argName := w.exprName(call.Args[0])
	return argName == sliceName
}

// containsPos checks if a node's range contains the given position.
func containsPos(node ast.Node, pos token.Pos) bool {
	return node.Pos() <= pos && pos <= node.End()
}
