package no_slice_bounds_out_of_range

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/lint"
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

// rangeScope tracks a for-range statement whose key variable is a safe index.
type rangeScope struct {
	keyName   string // the name of the range key variable (e.g. "i")
	sliceName string // the name of the ranged slice (e.g. "pes")
}

// sortSliceScope tracks a sort.Slice callback whose parameters are safe indices.
type sortSliceScope struct {
	paramNames []string // callback parameter names (e.g. ["i", "j"])
	sliceName  string   // the slice being sorted
}

func (w *lintNoSliceBoundsOutOfRange) checkFunctionBody(body *ast.BlockStmt) {
	w.walkBlock(body, body, nil, nil)
}

// walkBlock recursively walks AST nodes, tracking range and sort.Slice scopes.
func (w *lintNoSliceBoundsOutOfRange) walkBlock(node ast.Node, body *ast.BlockStmt, ranges []rangeScope, sortScopes []sortSliceScope) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch expr := n.(type) {
		case *ast.RangeStmt:
			w.walkRange(expr, body, ranges, sortScopes)
			return false // we handle children ourselves
		case *ast.CallExpr:
			if sliceName, paramNames := w.parseSortSliceCall(expr); sliceName != "" {
				newScope := sortSliceScope{paramNames: paramNames, sliceName: sliceName}
				w.walkSortSliceCallback(expr, body, ranges, append(sortScopes, newScope))
				return false
			}
		case *ast.IndexExpr:
			w.checkIndexExprWithScopes(expr, body, ranges, sortScopes)
		case *ast.SliceExpr:
			w.checkSliceExprWithScopes(expr, body, ranges, sortScopes)
		}
		return true
	})
}

// walkRange walks a for-range statement, adding its key variable as a safe index.
func (w *lintNoSliceBoundsOutOfRange) walkRange(rs *ast.RangeStmt, body *ast.BlockStmt, ranges []rangeScope, sortScopes []sortSliceScope) {
	// Only track if the range key is an identifier and the ranged expression is a named slice
	if rs.Key != nil && rs.Body != nil {
		keyIdent, ok := rs.Key.(*ast.Ident)
		if ok && keyIdent.Name != "_" && w.isSliceType(rs.X) {
			sliceName := w.exprName(rs.X)
			if sliceName != "" {
				newRanges := append(ranges, rangeScope{keyName: keyIdent.Name, sliceName: sliceName})
				w.walkBlock(rs.Body, body, newRanges, sortScopes)
				return
			}
		}
	}
	// Fallback: walk children with current scopes
	w.walkBlock(rs.Body, body, ranges, sortScopes)
}

// parseSortSliceCall checks if a call is sort.Slice(slice, func(i, j int) bool { ... })
// and returns the slice name and callback parameter names.
func (w *lintNoSliceBoundsOutOfRange) parseSortSliceCall(call *ast.CallExpr) (string, []string) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "sort" {
		return "", nil
	}
	if sel.Sel.Name != "Slice" && sel.Sel.Name != "SliceStable" {
		return "", nil
	}
	if len(call.Args) != 2 {
		return "", nil
	}

	sliceName := w.exprName(call.Args[0])
	if sliceName == "" {
		return "", nil
	}

	funcLit, ok := call.Args[1].(*ast.FuncLit)
	if !ok || funcLit.Type.Params == nil {
		return "", nil
	}

	var paramNames []string
	for _, field := range funcLit.Type.Params.List {
		for _, name := range field.Names {
			paramNames = append(paramNames, name.Name)
		}
	}
	if len(paramNames) == 0 {
		return "", nil
	}

	return sliceName, paramNames
}

// walkSortSliceCallback walks the sort.Slice call but avoids re-walking the callback body
// with different scopes. Instead, we walk the callback body with the new sort scope.
func (w *lintNoSliceBoundsOutOfRange) walkSortSliceCallback(call *ast.CallExpr, body *ast.BlockStmt, ranges []rangeScope, sortScopes []sortSliceScope) {
	funcLit, ok := call.Args[1].(*ast.FuncLit)
	if !ok || funcLit.Body == nil {
		return
	}
	// Walk the callback body with the added sort scope
	w.walkBlock(funcLit.Body, body, ranges, sortScopes)
}

func (w *lintNoSliceBoundsOutOfRange) isIndexSafe(indexExpr ast.Expr, sliceName string, ranges []rangeScope, sortScopes []sortSliceScope) bool {
	idxIdent, ok := indexExpr.(*ast.Ident)
	if !ok {
		return false
	}
	idxName := idxIdent.Name

	// Check range scopes
	for _, rs := range ranges {
		if rs.keyName == idxName && rs.sliceName == sliceName {
			return true
		}
	}

	// Check sort.Slice scopes
	for _, ss := range sortScopes {
		if ss.sliceName == sliceName {
			for _, pn := range ss.paramNames {
				if pn == idxName {
					return true
				}
			}
		}
	}

	return false
}

func (w *lintNoSliceBoundsOutOfRange) checkIndexExprWithScopes(expr *ast.IndexExpr, body *ast.BlockStmt, ranges []rangeScope, sortScopes []sortSliceScope) {
	if !w.isSliceType(expr.X) {
		return
	}

	sliceName := w.exprName(expr.X)
	if sliceName == "" {
		return
	}

	if w.isIndexSafe(expr.Index, sliceName, ranges, sortScopes) {
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

func (w *lintNoSliceBoundsOutOfRange) checkSliceExprWithScopes(expr *ast.SliceExpr, body *ast.BlockStmt, ranges []rangeScope, sortScopes []sortSliceScope) {
	if !w.isSliceType(expr.X) {
		return
	}

	sliceName := w.exprName(expr.X)
	if sliceName == "" {
		return
	}

	if expr.High == nil && expr.Low == nil {
		return // s[:] is always safe
	}

	// Check if low/high bounds are safe via range or sort scopes
	lowSafe := expr.Low == nil || w.isIndexSafe(expr.Low, sliceName, ranges, sortScopes)
	highSafe := expr.High == nil || w.isIndexSafe(expr.High, sliceName, ranges, sortScopes)
	if lowSafe && highSafe {
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
