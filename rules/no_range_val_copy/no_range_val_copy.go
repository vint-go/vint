package no_range_val_copy

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRangeValCopyRule detects range loops that copy large values on each iteration.
// When a range loop captures values (e.g., for _, x := range xs), each iteration
// copies the value. For large structs this copying is expensive.
type NoRangeValCopyRule struct {
	sizeThreshold int64
	skipTestFuncs bool
}

const defaultSizeThreshold = 128

// Configure validates and applies the rule configuration.
func (r *NoRangeValCopyRule) Configure(arguments lint.Arguments) error {
	r.sizeThreshold = defaultSizeThreshold
	r.skipTestFuncs = true

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		threshold, ok := lint.ToInt64(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "noRangeValCopy" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.sizeThreshold = threshold
		return nil
	}

	for k, v := range argKV {
		switch normalizeOption(k) {
		case "sizethreshold":
			threshold, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for sizeThreshold in "noRangeValCopy" rule; need integer but got %T`, v)
			}
			r.sizeThreshold = threshold
		case "skiptestfuncs":
			skip, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skipTestFuncs in "noRangeValCopy" rule; need bool but got %T`, v)
			}
			r.skipTestFuncs = skip
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoRangeValCopyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

	w := &lintRangeValCopy{
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
func (r *NoRangeValCopyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
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

	w := &lintRangeValCopy{
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
func (*NoRangeValCopyRule) Name() string {
	return "noRangeValCopy"
}

// Group returns the rule group.
func (*NoRangeValCopyRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoRangeValCopyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoRangeValCopyRule) RequiresTypecheck() bool {
	return true
}

type lintRangeValCopy struct {
	typesInfo     *types.Info
	sizes         types.Sizes
	sizeThreshold int64
	onFailure     func(lint.Failure)
}

func (w *lintRangeValCopy) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// We only care about the value variable (second variable in range).
	// for _, val := range xs { ... }
	if rangeStmt.Value == nil {
		return w
	}

	// The value must be an identifier (not blank).
	valIdent, ok := rangeStmt.Value.(*ast.Ident)
	if !ok || valIdent.Name == "_" {
		return w
	}

	// Get the type of the value variable.
	valType := w.typesInfo.TypeOf(rangeStmt.Value)
	if valType == nil {
		return w
	}

	// Skip pointer types, interfaces, and other reference types.
	if isReferenceType(valType) {
		return w
	}

	size, ok := safeSizeof(w.sizes, valType)
	if !ok {
		return w
	}

	if size >= w.sizeThreshold {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryOptimization,
			Confidence: 1,
			Node:       rangeStmt,
			Failure:    fmt.Sprintf("range value '%s' copies %d bytes each iteration, consider using index access or taking the address", valIdent.Name, size),
		})
	}

	return w
}

// isReferenceType returns true if the type is a reference type (pointer, slice,
// map, channel, interface, function) that does not involve expensive copying.
func isReferenceType(t types.Type) bool {
	switch t.Underlying().(type) {
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Interface, *types.Signature:
		return true
	}
	return false
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
