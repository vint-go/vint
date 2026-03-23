package no_excessive_shift

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExcessiveShiftRule checks for shifts that exceed the width of an integer.
type NoExcessiveShiftRule struct{}

// Apply applies the rule to given file.
func (r *NoExcessiveShiftRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintExcessiveShift{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoExcessiveShiftRule) Name() string {
	return "noExcessiveShift"
}

// Group returns the rule group.
func (*NoExcessiveShiftRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoExcessiveShiftRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoExcessiveShiftRule) RequiresTypecheck() bool {
	return true
}

type lintExcessiveShift struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintExcessiveShift) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.SHL && binExpr.Op != token.SHR {
		return w
	}

	// Determine the bit width of the left operand's type.
	typ := w.file.Pkg.TypeOf(binExpr.X)
	if typ == nil {
		return w
	}

	bitWidth := intTypeBitWidth(typ.Underlying())
	if bitWidth == 0 {
		return w
	}

	// Determine the shift amount from the right operand.
	// We need it to be a compile-time constant.
	shiftAmount, ok := w.constShiftAmount(binExpr.Y)
	if !ok {
		return w
	}

	if shiftAmount >= bitWidth {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       binExpr,
			Failure:    fmt.Sprintf("shift count %d exceeds the bit width of type %s (%d bits)", shiftAmount, typ.String(), bitWidth),
		})
	}

	return w
}

// constShiftAmount tries to extract a compile-time constant integer shift amount.
func (w *lintExcessiveShift) constShiftAmount(expr ast.Expr) (int64, bool) {
	info := w.file.Pkg.TypesInfo()
	if info == nil {
		return 0, false
	}

	tv, ok := info.Types[expr]
	if !ok || tv.Value == nil {
		return 0, false
	}

	if tv.Value.Kind() != constant.Int {
		return 0, false
	}

	val, exact := constant.Int64Val(tv.Value)
	if !exact {
		return 0, false
	}

	return val, true
}

// intTypeBitWidth returns the bit width of an integer type, or 0 if it is not
// a fixed-width integer type.
func intTypeBitWidth(t types.Type) int64 {
	basic, ok := t.(*types.Basic)
	if !ok {
		return 0
	}

	switch basic.Kind() {
	case types.Int8, types.Uint8:
		return 8
	case types.Int16, types.Uint16:
		return 16
	case types.Int32, types.Uint32:
		return 32
	case types.Int64, types.Uint64:
		return 64
	case types.Int, types.Uint, types.Uintptr:
		// Platform-dependent size; assume 64-bit as the common case.
		// We only flag when shift >= bit width, so using 64 avoids false positives
		// on 64-bit platforms while still catching clearly excessive shifts on 32-bit.
		return 64
	default:
		return 0
	}
}
