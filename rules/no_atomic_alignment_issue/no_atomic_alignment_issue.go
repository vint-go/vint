package no_atomic_alignment_issue

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoAtomicAlignmentIssueRule checks for non-64-bit-aligned arguments to
// sync/atomic functions. On 32-bit platforms, 64-bit atomic operations
// require that the variable be 64-bit aligned in memory.
type NoAtomicAlignmentIssueRule struct{}

// Apply applies the rule to given file.
func (r *NoAtomicAlignmentIssueRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoAtomicAlignment{
		typesInfo: file.Pkg.TypesInfo(),
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoAtomicAlignmentIssueRule) Name() string {
	return "noAtomicAlignmentIssue"
}

// Group returns the rule group.
func (*NoAtomicAlignmentIssueRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoAtomicAlignmentIssueRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoAtomicAlignmentIssueRule) RequiresTypecheck() bool {
	return true
}

// atomic64BitFunctions is the set of sync/atomic functions that operate on
// 64-bit values and therefore require 64-bit alignment on 32-bit platforms.
var atomic64BitFunctions = map[string]bool{
	"AddInt64":             true,
	"LoadInt64":            true,
	"StoreInt64":           true,
	"CompareAndSwapInt64":  true,
	"SwapInt64":            true,
	"AddUint64":            true,
	"LoadUint64":           true,
	"StoreUint64":          true,
	"CompareAndSwapUint64": true,
	"SwapUint64":           true,
}

type lintNoAtomicAlignment struct {
	typesInfo *types.Info
	onFailure func(lint.Failure)
}

func (w *lintNoAtomicAlignment) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if !atomic64BitFunctions[sel.Sel.Name] {
		return w
	}

	// Verify the call is from sync/atomic using type info.
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}
	if w.typesInfo != nil {
		pkgName, ok := w.typesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || pkgName.Imported().Path() != "sync/atomic" {
			return w
		}
	}

	if len(call.Args) == 0 {
		return w
	}

	// The first argument should be a pointer to the 64-bit variable.
	// We're looking for &structExpr.field patterns.
	arg := call.Args[0]
	uarg, ok := arg.(*ast.UnaryExpr)
	if !ok || uarg.Op != token.AND {
		return w
	}

	// Check if the operand is a field selector (e.g., c.count).
	fieldSel, ok := uarg.X.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if w.typesInfo == nil {
		return w
	}

	// Get the type of the struct expression (e.g., type of 'c' in 'c.count').
	structType := w.typesInfo.TypeOf(fieldSel.X)
	if structType == nil {
		return w
	}

	// Dereference pointer if needed (e.g., *Counter -> Counter).
	if ptr, ok := structType.(*types.Pointer); ok {
		structType = ptr.Elem()
	}

	// Get the underlying struct type.
	st, ok := structType.Underlying().(*types.Struct)
	if !ok {
		return w
	}

	// Find the field and check its offset.
	fieldName := fieldSel.Sel.Name
	fieldIndex := -1
	for i := 0; i < st.NumFields(); i++ {
		if st.Field(i).Name() == fieldName {
			fieldIndex = i
			break
		}
	}

	if fieldIndex < 0 {
		return w
	}

	// Use a 32-bit sizes model to compute field offsets.
	// On 32-bit platforms, the maximum alignment is 4 bytes.
	sizes := &types.StdSizes{WordSize: 4, MaxAlign: 4}
	offsets := sizes.Offsetsof(structFields(st))
	if fieldIndex >= len(offsets) {
		return w
	}

	offset := offsets[fieldIndex]
	if offset%8 != 0 {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "address of non 64-bit aligned field ." + fieldName + " passed to atomic function",
		})
	}

	return w
}

// structFields extracts the field list from a struct type as a slice of *types.Var.
func structFields(st *types.Struct) []*types.Var {
	fields := make([]*types.Var, st.NumFields())
	for i := 0; i < st.NumFields(); i++ {
		fields[i] = st.Field(i)
	}
	return fields
}
