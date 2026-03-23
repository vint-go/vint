package no_suspicious_sort_slice

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSuspiciousSortSliceRule detects suspicious sort.Slice and sort.SliceStable
// calls where the comparison function either references a different slice than
// the one being sorted, or uses reversed index parameters.
type NoSuspiciousSortSliceRule struct{}

// Apply applies the rule to given file.
func (r *NoSuspiciousSortSliceRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousSortSlice{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSuspiciousSortSliceRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousSortSlice{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSuspiciousSortSliceRule) Name() string {
	return "noSuspiciousSortSlice"
}

// Group returns the rule group.
func (*NoSuspiciousSortSliceRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSuspiciousSortSliceRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// sortSliceFuncNames lists the sort package functions to check.
var sortSliceFuncNames = map[string]bool{
	"Slice":       true,
	"SliceStable": true,
}

type lintSuspiciousSortSlice struct {
	onFailure func(lint.Failure)
}

func (w *lintSuspiciousSortSlice) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for sort.Slice or sort.SliceStable
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if !sortSliceFuncNames[sel.Sel.Name] {
		return w
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "sort" {
		return w
	}

	// Need exactly 2 args: the slice and the less function
	if len(call.Args) != 2 {
		return w
	}

	sliceArg := call.Args[0]
	lessFunc := call.Args[1]

	// The less function must be a func literal
	funcLit, ok := lessFunc.(*ast.FuncLit)
	if !ok {
		return w
	}

	// The less function must have exactly 2 parameters (i, j int)
	if funcLit.Type.Params == nil || len(funcLit.Type.Params.List) == 0 {
		return w
	}

	// Extract parameter names for the two index params
	var paramI, paramJ string
	params := funcLit.Type.Params.List
	var allNames []*ast.Ident
	for _, p := range params {
		allNames = append(allNames, p.Names...)
	}
	if len(allNames) != 2 {
		return w
	}
	paramI = allNames[0].Name
	paramJ = allNames[1].Name

	// Get the name of the slice being sorted
	sliceName := astutils.GoFmt(sliceArg)

	// Collect all index expressions in the function body
	indexExprs := collectIndexExprs(funcLit.Body)
	if len(indexExprs) == 0 {
		return w
	}

	// Check pattern 1: missing slice reference
	// The comparison function indexes into a different collection than the slice being sorted
	w.checkMissingSliceRef(call, sliceName, indexExprs, paramI, paramJ)

	// Check pattern 2: reversed index parameters
	// The comparison uses j on left and i on right, which may indicate reversed sorting logic
	w.checkReversedIndices(call, sliceName, indexExprs, paramI, paramJ)

	return w
}

// checkMissingSliceRef reports if the comparison function never references the
// slice being sorted via indexing with the parameter names.
func (w *lintSuspiciousSortSlice) checkMissingSliceRef(
	call *ast.CallExpr,
	sliceName string,
	indexExprs []*ast.IndexExpr,
	paramI, paramJ string,
) {
	for _, idx := range indexExprs {
		collName := astutils.GoFmt(idx.X)
		idxName := identName(idx.Index)
		if collName == sliceName && (idxName == paramI || idxName == paramJ) {
			return // found a valid reference to the sorted slice
		}
	}

	// None of the index expressions reference the sorted slice
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryLogic,
		Failure:    fmt.Sprintf("sort.Slice's comparison function does not reference the sorted slice %s", sliceName),
	})
}

// checkReversedIndices reports if the comparison uses j before i (reversed),
// specifically checking for patterns like slice[j] < slice[i].
func (w *lintSuspiciousSortSlice) checkReversedIndices(
	call *ast.CallExpr,
	sliceName string,
	indexExprs []*ast.IndexExpr,
	paramI, paramJ string,
) {
	// Look for binary expressions comparing slice[j] op slice[i]
	binaryExprs := collectBinaryExprs(call.Args[1])
	for _, binExpr := range binaryExprs {
		leftIdx, leftOk := binExpr.X.(*ast.IndexExpr)
		rightIdx, rightOk := binExpr.Y.(*ast.IndexExpr)
		if !leftOk || !rightOk {
			continue
		}

		leftColl := astutils.GoFmt(leftIdx.X)
		rightColl := astutils.GoFmt(rightIdx.X)
		leftIndex := identName(leftIdx.Index)
		rightIndex := identName(rightIdx.Index)

		// Both sides must reference the same sorted slice
		if leftColl != sliceName || rightColl != sliceName {
			continue
		}

		// Check if indices are reversed: j on left, i on right
		if leftIndex == paramJ && rightIndex == paramI {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("sort.Slice's comparison function has reversed indices: %s[%s] compared with %s[%s]", sliceName, paramJ, sliceName, paramI),
			})
			return // report once per call
		}
	}
}

// identName returns the name of an identifier expression, or empty string.
func identName(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

// collectIndexExprs collects all ast.IndexExpr nodes in the given AST subtree.
func collectIndexExprs(node ast.Node) []*ast.IndexExpr {
	var result []*ast.IndexExpr
	ast.Inspect(node, func(n ast.Node) bool {
		if idx, ok := n.(*ast.IndexExpr); ok {
			result = append(result, idx)
		}
		return true
	})
	return result
}

// collectBinaryExprs collects all ast.BinaryExpr nodes in the given AST subtree.
func collectBinaryExprs(node ast.Node) []*ast.BinaryExpr {
	var result []*ast.BinaryExpr
	ast.Inspect(node, func(n ast.Node) bool {
		if bin, ok := n.(*ast.BinaryExpr); ok {
			result = append(result, bin)
		}
		return true
	})
	return result
}
