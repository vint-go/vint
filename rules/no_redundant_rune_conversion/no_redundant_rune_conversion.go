package no_redundant_rune_conversion

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantRuneConversionRule detects unnecessary conversions of a string
// to []rune before ranging over it. The for-range loop already iterates
// over runes natively, so the conversion only adds an allocation.
type NoRedundantRuneConversionRule struct{}

// Apply applies the rule to the given file.
func (r *NoRedundantRuneConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantRuneConversion{
		typesInfo: typesInfo,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantRuneConversionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantRuneConversion{
		typesInfo: typesInfo,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantRuneConversionRule) Name() string {
	return "noRedundantRuneConversion"
}

// Group returns the rule group.
func (*NoRedundantRuneConversionRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantRuneConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoRedundantRuneConversionRule) RequiresTypecheck() bool {
	return true
}

type lintRedundantRuneConversion struct {
	typesInfo *types.Info
	onFailure func(lint.Failure)
}

func (w *lintRedundantRuneConversion) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Check if the range expression is a type conversion call: []rune(expr)
	call, ok := rangeStmt.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	// The "function" in a type conversion like []rune(s) is an *ast.ArrayType
	arrType, ok := call.Fun.(*ast.ArrayType)
	if !ok {
		return w
	}

	// Must be a slice (Len == nil), not an array
	if arrType.Len != nil {
		return w
	}

	// The element type must be "rune"
	eltIdent, ok := arrType.Elt.(*ast.Ident)
	if !ok || eltIdent.Name != "rune" {
		return w
	}

	// Must have exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	// Use type info to confirm the argument is a string
	arg := call.Args[0]
	argType := w.typesInfo.TypeOf(arg)
	if argType == nil {
		return w
	}

	// Check that the underlying type is string
	basic, ok := argType.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryOptimization,
		Confidence: 1,
		Node:       call,
		Failure:    "unnecessary conversion to []rune before range loop",
	})

	return w
}
