package no_huge_param

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHugeParamRule detects function parameters that exceed a specified byte
// threshold, as passing large structures by value causes unnecessary copying
// overhead. The String() string method is automatically excluded to avoid
// flagging Stringer interface implementations.
type NoHugeParamRule struct {
	sizeThreshold int64
}

const defaultSizeThreshold = 80

// Configure validates and applies the rule configuration.
func (r *NoHugeParamRule) Configure(arguments lint.Arguments) error {
	r.sizeThreshold = defaultSizeThreshold

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// Try direct int64 argument.
		threshold, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noHugeParam" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.sizeThreshold = threshold
		return nil
	}

	for k, v := range argKV {
		if normalizeOption(k) == "sizethreshold" {
			threshold, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for sizeThreshold in "noHugeParam" rule; need int64 but got %T`, v)
			}
			r.sizeThreshold = threshold
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoHugeParamRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.sizeThreshold <= 0 {
		r.sizeThreshold = defaultSizeThreshold
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	sizes := types.SizesFor("gc", "amd64")

	var failures []lint.Failure

	w := &lintHugeParam{
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
func (r *NoHugeParamRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if r.sizeThreshold <= 0 {
		r.sizeThreshold = defaultSizeThreshold
	}

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	sizes := types.SizesFor("gc", "amd64")

	var failures []lint.Failure

	w := &lintHugeParam{
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
func (*NoHugeParamRule) Name() string {
	return "noHugeParam"
}

// Group returns the rule group.
func (*NoHugeParamRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoHugeParamRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoHugeParamRule) RequiresTypecheck() bool {
	return true
}

type lintHugeParam struct {
	typesInfo     *types.Info
	sizes         types.Sizes
	sizeThreshold int64
	onFailure     func(lint.Failure)
}

func (w *lintHugeParam) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	// Exclude String() string methods (Stringer interface).
	if isStringerMethod(funcDecl) {
		return w
	}

	// Check method receiver for huge value receivers.
	if funcDecl.Recv != nil {
		for _, recv := range funcDecl.Recv.List {
			w.checkField(recv)
		}
	}

	// Check function parameters.
	if funcDecl.Type != nil && funcDecl.Type.Params != nil {
		for _, param := range funcDecl.Type.Params.List {
			w.checkField(param)
		}
	}

	return w
}

// checkField checks a single field (parameter or receiver) for exceeding the size threshold.
func (w *lintHugeParam) checkField(field *ast.Field) {
	// Skip pointer types -- already passed by reference.
	if _, ok := field.Type.(*ast.StarExpr); ok {
		return
	}

	paramType := w.typesInfo.TypeOf(field.Type)
	if paramType == nil {
		return
	}

	size, ok := safeSizeof(w.sizes, paramType)
	if !ok {
		return
	}

	if size > w.sizeThreshold {
		if len(field.Names) == 0 {
			// Unnamed parameter.
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryOptimization,
				Confidence: 1,
				Node:       field,
				Failure:    fmt.Sprintf("parameter exceeds the size threshold of %d bytes with a size of %d bytes, consider passing it by pointer", w.sizeThreshold, size),
			})
		} else {
			for _, name := range field.Names {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryOptimization,
					Confidence: 1,
					Node:       field,
					Failure:    fmt.Sprintf("parameter '%s' exceeds the size threshold of %d bytes with a size of %d bytes, consider passing it by pointer", name.Name, w.sizeThreshold, size),
				})
			}
		}
	}
}

// isStringerMethod returns true if the function declaration is a method
// named "String" that returns a single string result and has no parameters
// (matching the fmt.Stringer interface).
func isStringerMethod(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Recv == nil || funcDecl.Recv.NumFields() == 0 {
		return false
	}
	if funcDecl.Name.Name != "String" {
		return false
	}
	// Must have no parameters.
	if funcDecl.Type.Params != nil && funcDecl.Type.Params.NumFields() > 0 {
		return false
	}
	// Must return exactly one result of type string.
	if funcDecl.Type.Results == nil || funcDecl.Type.Results.NumFields() != 1 {
		return false
	}
	retType, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
	if !ok {
		return false
	}
	return retType.Name == "string"
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
