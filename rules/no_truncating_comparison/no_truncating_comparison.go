package no_truncating_comparison

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoTruncatingComparisonRule detects potential truncation issues when comparing
// ints of different sizes. It flags comparisons where an integer is cast to a
// smaller type before being compared, since casting the narrower operand to the
// larger type would be safer and more correct.
type NoTruncatingComparisonRule struct {
	skipArchDependent bool
}

// Configure validates and applies the rule configuration.
func (r *NoTruncatingComparisonRule) Configure(arguments lint.Arguments) error {
	r.skipArchDependent = true // default

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noTruncatingComparison" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == "skiparchdependent" {
			skip, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skipArchDependent in "noTruncatingComparison" rule; need bool but got %T`, v)
			}
			r.skipArchDependent = skip
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoTruncatingComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	w := &lintTruncatingComparison{
		typesInfo:         typesInfo,
		skipArchDependent: r.skipArchDependent,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTruncatingComparisonRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	w := &lintTruncatingComparison{
		typesInfo:         typesInfo,
		skipArchDependent: r.skipArchDependent,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoTruncatingComparisonRule) Name() string {
	return "noTruncatingComparison"
}

// Group returns the rule group.
func (*NoTruncatingComparisonRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoTruncatingComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoTruncatingComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintTruncatingComparison struct {
	typesInfo         *types.Info
	skipArchDependent bool
	onFailure         func(lint.Failure)
}

func (w *lintTruncatingComparison) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if !isComparisonOp(binExpr.Op.String()) {
		return w
	}

	// Check if either side has a truncating conversion.
	w.checkSide(binExpr, binExpr.X, binExpr.Y)
	w.checkSide(binExpr, binExpr.Y, binExpr.X)

	return w
}

// checkSide checks if 'side' is a type conversion to a smaller int type,
// where the converted value's original type is larger than the target type.
// 'other' is the other side of the comparison.
func (w *lintTruncatingComparison) checkSide(binExpr *ast.BinaryExpr, side, other ast.Expr) {
	callExpr, ok := side.(*ast.CallExpr)
	if !ok {
		return
	}

	// Must be a type conversion (not a function call).
	// Type conversions have exactly one argument.
	if len(callExpr.Args) != 1 {
		return
	}

	// The function part must be an identifier that names a type.
	convType := w.typesInfo.TypeOf(callExpr.Fun)
	if convType == nil {
		return
	}

	// Resolve the target type of the conversion.
	targetType, ok := resolveBasicInt(convType)
	if !ok {
		return
	}

	// Get the type of the inner argument (the value being converted).
	innerArg := callExpr.Args[0]
	innerType := w.typesInfo.TypeOf(innerArg)
	if innerType == nil {
		return
	}

	innerBasic, ok := resolveBasicInt(innerType)
	if !ok {
		return
	}

	// Skip architecture-dependent types if configured.
	if w.skipArchDependent {
		if isArchDependent(targetType) || isArchDependent(innerBasic) {
			return
		}
	}

	targetSize := intBitSize(targetType)
	innerSize := intBitSize(innerBasic)

	if targetSize == 0 || innerSize == 0 {
		return
	}

	// Only flag if the conversion narrows (target is smaller than inner).
	if targetSize >= innerSize {
		return
	}

	// Check that the signs match for a fair comparison of sizes.
	// Signed-to-unsigned or vice versa conversions with narrowing are also suspicious.
	targetName := targetType.Name()
	innerName := innerBasic.Name()

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 0.8,
		Node:       binExpr,
		Failure:    fmt.Sprintf("truncating conversion from %s to %s before comparison may lose information", innerName, targetName),
	})
}

// resolveBasicInt returns the underlying *types.Basic if it is an integer type.
func resolveBasicInt(t types.Type) (*types.Basic, bool) {
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return nil, false
	}
	if basic.Info()&types.IsInteger == 0 {
		return nil, false
	}
	return basic, true
}

// isArchDependent returns true for int, uint, and uintptr types whose size
// varies by architecture.
func isArchDependent(b *types.Basic) bool {
	switch b.Kind() {
	case types.Int, types.Uint, types.Uintptr:
		return true
	}
	return false
}

// intBitSize returns the bit size of a basic integer type.
// Returns 0 for architecture-dependent types (int, uint, uintptr).
func intBitSize(b *types.Basic) int {
	switch b.Kind() {
	case types.Int8, types.Uint8:
		return 8
	case types.Int16, types.Uint16:
		return 16
	case types.Int32, types.Uint32:
		return 32
	case types.Int64, types.Uint64:
		return 64
	case types.Int, types.Uint, types.Uintptr:
		// Architecture-dependent; assume 64-bit for detection purposes.
		return 64
	}
	return 0
}

// isComparisonOp returns true if the operator string is a comparison operator.
func isComparisonOp(op string) bool {
	switch op {
	case "<", ">", "<=", ">=", "==", "!=":
		return true
	}
	return false
}

// normalizeOption normalizes a configuration option name by removing hyphens,
// underscores, and lowering case.
func normalizeOption(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}
