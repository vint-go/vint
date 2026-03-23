package use_optimal_field_alignment

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseOptimalFieldAlignmentRule detects structs that would use less memory if their fields were sorted.
type UseOptimalFieldAlignmentRule struct{}

// Apply applies the rule to given file.
func (r *UseOptimalFieldAlignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	// Use standard sizes for amd64 (8-byte word, 8-byte max alignment).
	sizes := types.SizesFor("gc", "amd64")

	w := &lintFieldAlignment{
		typesInfo: typesInfo,
		sizes:     sizes,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseOptimalFieldAlignmentRule) Name() string {
	return "useOptimalFieldAlignment"
}

// Group returns the rule group.
func (*UseOptimalFieldAlignmentRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseOptimalFieldAlignmentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*UseOptimalFieldAlignmentRule) RequiresTypecheck() bool {
	return true
}

type fieldInfo struct {
	align int64
	size  int64
	name  string
}

type lintFieldAlignment struct {
	typesInfo *types.Info
	sizes     types.Sizes
	onFailure func(lint.Failure)
}

func (w *lintFieldAlignment) Visit(node ast.Node) ast.Visitor {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return w
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return w
	}

	if structType.Fields == nil || structType.Fields.NumFields() < 2 {
		return w
	}

	// Get the type object for this type spec.
	obj := w.typesInfo.Defs[typeSpec.Name]
	if obj == nil {
		return w
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return w
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return w
	}

	numFields := st.NumFields()
	if numFields < 2 {
		return w
	}

	currentSize, ok := safeSizeof(w.sizes, st)
	if !ok {
		return w
	}

	// Calculate optimal size by sorting fields by alignment (descending),
	// then by size (descending) for fields with equal alignment.
	fields := make([]fieldInfo, numFields)
	for i := 0; i < numFields; i++ {
		f := st.Field(i)
		ft := f.Type()
		a, aOk := safeAlignof(w.sizes, ft)
		s, sOk := safeSizeof(w.sizes, ft)
		if !aOk || !sOk {
			return w
		}
		fields[i] = fieldInfo{
			align: a,
			size:  s,
			name:  f.Name(),
		}
	}

	// Sort by alignment descending, then by size descending.
	sort.SliceStable(fields, func(i, j int) bool {
		if fields[i].align != fields[j].align {
			return fields[i].align > fields[j].align
		}
		return fields[i].size > fields[j].size
	})

	// Calculate optimized size using the sorted field order.
	optimalSize := calculateStructSize(fields)

	if optimalSize < currentSize {
		optimalOrder := make([]string, len(fields))
		for i, f := range fields {
			optimalOrder[i] = f.name
		}

		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryOptimization,
			Confidence: 1,
			Node:       typeSpec,
			Failure: fmt.Sprintf("struct %s could have its fields rearranged to use %d bytes instead of %d bytes (optimal field order: %s)",
				typeSpec.Name.Name, optimalSize, currentSize, strings.Join(optimalOrder, ", ")),
		})
	}

	return w
}

// calculateStructSize calculates the total size of a struct given ordered fields.
func calculateStructSize(fields []fieldInfo) int64 {
	if len(fields) == 0 {
		return 0
	}

	var offset int64
	var maxAlign int64

	for _, f := range fields {
		a := f.align
		if a > maxAlign {
			maxAlign = a
		}
		// Align offset to field alignment.
		offset = align(offset, a)
		offset += f.size
	}

	// Add trailing padding to align struct to its overall alignment.
	if maxAlign > 0 {
		offset = align(offset, maxAlign)
	}

	return offset
}

// align rounds up n to the next multiple of a.
func align(n, a int64) int64 {
	return (n + a - 1) &^ (a - 1)
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

// safeAlignof wraps types.Sizes.Alignof, recovering from panics.
func safeAlignof(sizes types.Sizes, t types.Type) (size int64, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	return sizes.Alignof(t), true
}
