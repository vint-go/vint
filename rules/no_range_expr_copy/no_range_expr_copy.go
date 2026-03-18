package no_range_expr_copy

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRangeExprCopyRule detects expensive copies of for-range loop range
// expressions. When iterating over an array in a for-range loop, Go copies
// the entire array. This rule warns when the copied data size reaches the
// configured threshold, suggesting the use of a pointer to the array (with &)
// to avoid the unnecessary copy.
type NoRangeExprCopyRule struct {
	sizeThreshold int64
	skipTestFuncs bool
}

const defaultSizeThreshold = 512

// Configure validates and applies the rule configuration.
func (r *NoRangeExprCopyRule) Configure(arguments lint.Arguments) error {
	r.sizeThreshold = defaultSizeThreshold
	r.skipTestFuncs = true

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		threshold, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noRangeExprCopy" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.sizeThreshold = threshold
		return nil
	}

	for k, v := range argKV {
		switch normalizeOption(k) {
		case "sizethreshold":
			threshold, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for sizeThreshold in "noRangeExprCopy" rule; need int64 but got %T`, v)
			}
			r.sizeThreshold = threshold
		case "skiptestfuncs":
			skip, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skipTestFuncs in "noRangeExprCopy" rule; need bool but got %T`, v)
			}
			r.skipTestFuncs = skip
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoRangeExprCopyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.sizeThreshold <= 0 {
		r.sizeThreshold = defaultSizeThreshold
	}

	if r.skipTestFuncs && file.IsTest() {
		return nil
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	sizes := types.SizesFor("gc", "amd64")

	var failures []lint.Failure

	w := &lintRangeExprCopy{
		typesInfo:     typesInfo,
		sizes:         sizes,
		sizeThreshold: r.sizeThreshold,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRangeExprCopyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if r.sizeThreshold <= 0 {
		r.sizeThreshold = defaultSizeThreshold
	}

	if r.skipTestFuncs && file.IsTest() {
		return nil
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	sizes := types.SizesFor("gc", "amd64")

	var failures []lint.Failure

	w := &lintRangeExprCopy{
		typesInfo:     typesInfo,
		sizes:         sizes,
		sizeThreshold: r.sizeThreshold,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRangeExprCopyRule) Name() string {
	return "noRangeExprCopy"
}

// Group returns the rule group.
func (*NoRangeExprCopyRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoRangeExprCopyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoRangeExprCopyRule) RequiresTypecheck() bool {
	return true
}

type lintRangeExprCopy struct {
	typesInfo     *types.Info
	sizes         types.Sizes
	sizeThreshold int64
	onFailure     func(lint.Failure)
}

func (w *lintRangeExprCopy) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Get the type of the range expression.
	exprType := w.typesInfo.TypeOf(rangeStmt.X)
	if exprType == nil {
		return w
	}

	// Unwrap the underlying type to check if it's an array.
	underlying := exprType.Underlying()
	arrType, isArray := underlying.(*types.Array)
	if !isArray {
		return w
	}

	// Compute the size of the array.
	elemSize, ok := safeSizeof(w.sizes, arrType.Elem())
	if !ok {
		return w
	}
	totalSize := elemSize * arrType.Len()

	if totalSize < w.sizeThreshold {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryOptimization,
		Confidence: 1,
		Node:       rangeStmt,
		Failure:    fmt.Sprintf("range expression copies %d bytes, consider using a pointer (e.g. &expr)", totalSize),
	})

	return w
}

// safeSizeof wraps types.Sizes.Sizeof, recovering from panics caused by
// types that cannot be sized (e.g. generic type parameters).
func safeSizeof(sizes types.Sizes, t types.Type) (size int64, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	return sizes.Sizeof(t), true
}

// normalizeOption normalizes a configuration option name by removing hyphens,
// underscores, and lowering case.
func normalizeOption(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}
