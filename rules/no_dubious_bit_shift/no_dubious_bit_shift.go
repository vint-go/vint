package no_dubious_bit_shift

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDubiousBitShiftRule detects dubious bit shifting of fixed-size integer values.
// Shifting an integer by more bits than its size always results in zero (for unsigned)
// or zero/-1 (for signed).
type NoDubiousBitShiftRule struct{}

// Apply applies the rule to given file.
func (r *NoDubiousBitShiftRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintDubiousBitShift{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDubiousBitShiftRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDubiousBitShift{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDubiousBitShiftRule) Name() string {
	return "noDubiousBitShift"
}

// Group returns the rule group.
func (*NoDubiousBitShiftRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDubiousBitShiftRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*NoDubiousBitShiftRule) RequiresTypecheck() bool {
	return true
}

type lintDubiousBitShift struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintDubiousBitShift) Visit(node ast.Node) ast.Visitor {
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
			Failure:    fmt.Sprintf("dubious bit shift of %s value by %d bits, the type has only %d bits", typ.String(), shiftAmount, bitWidth),
		})
	}

	return w
}

// constShiftAmount tries to extract a compile-time constant integer shift amount.
func (w *lintDubiousBitShift) constShiftAmount(expr ast.Expr) (int64, bool) {
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
// a fixed-width integer type. Platform-dependent sizes (int, uint, uintptr)
// return 0 because their width varies by platform.
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
	default:
		// int, uint, uintptr have platform-dependent sizes;
		// skip them to avoid false positives.
		return 0
	}
}
