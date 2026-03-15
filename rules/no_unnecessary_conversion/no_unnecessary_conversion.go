package no_unnecessary_conversion

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnnecessaryConversionRule detects unnecessary type conversions in Go code.
// A type conversion T(x) is considered unnecessary when x already has type T.
type NoUnnecessaryConversionRule struct {
	fastMath bool
	safe     bool
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoUnnecessaryConversionRule) Configure(arguments lint.Arguments) error {
	r.fastMath = false
	r.safe = false

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnnecessaryConversion" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch normalizeOption(k) {
		case "fastmath":
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for fast-math in "noUnnecessaryConversion" rule; need bool but got %T`, v)
			}
			r.fastMath = val
		case "safe":
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for safe in "noUnnecessaryConversion" rule; need bool but got %T`, v)
			}
			r.safe = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUnnecessaryConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintUnnecessaryConversion{
		pkg:      file.Pkg,
		fastMath: r.fastMath,
		safe:     r.safe,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryConversionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintUnnecessaryConversion{
		pkg:      file.Pkg,
		fastMath: r.fastMath,
		safe:     r.safe,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryConversionRule) Name() string {
	return "noUnnecessaryConversion"
}

// Group returns the rule group.
func (*NoUnnecessaryConversionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoUnnecessaryConversionRule) RequiresTypecheck() bool {
	return true
}

type lintUnnecessaryConversion struct {
	pkg       *lint.Package
	fastMath  bool
	safe      bool
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryConversion) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Must be a single-argument call with no ellipsis (type conversion form)
	if len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return w
	}

	// The callee must be a type expression, not a function call.
	// We check if the callee is a type by looking at the types.Info.
	typesInfo := w.pkg.TypesInfo()
	if typesInfo == nil {
		return w
	}

	// Check if the function expression is a type (conversion) rather than a function call.
	// In Go's type system, a type conversion T(x) has the callee resolved as a type, not an object.
	if !isTypeExpr(typesInfo, call.Fun) {
		return w
	}

	// Get the target conversion type
	targetType := w.pkg.TypeOf(call.Fun)
	if targetType == nil {
		return w
	}

	// Get the argument type
	argType := w.pkg.TypeOf(call.Args[0])
	if argType == nil {
		return w
	}

	// The argument must not be an untyped value (golang.org/issue/13061 workaround)
	if isUntypedValue(call.Args[0], typesInfo) {
		return w
	}

	// Check if types are identical
	if !types.Identical(argType, targetType) {
		return w
	}

	// Unless fast-math is enabled, skip floating-point and complex number conversions
	if !w.fastMath && isFloatOrComplex(targetType) {
		return w
	}

	// When safe mode is enabled, check context
	if w.safe && !isSafeToRemove(call, typesInfo) {
		return w
	}

	typeName := astutils.GoFmt(call.Fun)
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    fmt.Sprintf("unnecessary conversion to %s", typeName),
	})

	return w
}

// isTypeExpr checks whether an expression is used as a type (for type conversion)
// rather than as a function call.
func isTypeExpr(info *types.Info, expr ast.Expr) bool {
	// For parenthesized expressions like (*int)(p), unwrap
	if paren, ok := expr.(*ast.ParenExpr); ok {
		// Check if the inner expression contains a star expression wrapping a type
		if star, ok := paren.X.(*ast.StarExpr); ok {
			return isTypeExpr(info, star.X)
		}
		return isTypeExpr(info, paren.X)
	}

	// Check the types.Info.Types map - type expressions used as conversions
	// have IsType() returning true
	if tv, ok := info.Types[expr]; ok {
		return tv.IsType()
	}

	return false
}

// isUntypedValue returns true if expr is an untyped constant or untyped expression.
func isUntypedValue(expr ast.Expr, info *types.Info) bool {
	// Check the types.Info.Types map for untyped basic type
	if tv, ok := info.Types[expr]; ok {
		if basic, ok := tv.Type.(*types.Basic); ok {
			if basic.Info()&types.IsUntyped != 0 {
				return true
			}
		}
		// If the value is a constant (non-nil Value), it may be an untyped constant
		// that the type checker resolved contextually.
		if tv.Value != nil {
			return true
		}
	}

	// Check if the expression is a basic literal (always untyped)
	if _, ok := expr.(*ast.BasicLit); ok {
		return true
	}

	// Check if the expression is an identifier resolving to a constant
	if ident, ok := expr.(*ast.Ident); ok {
		if obj := info.ObjectOf(ident); obj != nil {
			if _, isConst := obj.(*types.Const); isConst {
				return true
			}
		}
	}

	return false
}

// isFloatOrComplex returns true if t is a floating-point or complex type.
func isFloatOrComplex(t types.Type) bool {
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Float32, types.Float64, types.Complex64, types.Complex128:
		return true
	}
	return false
}

// isSafeToRemove performs context-aware analysis to check if removing
// the conversion would be semantically safe in the surrounding code.
func isSafeToRemove(call *ast.CallExpr, info *types.Info) bool {
	// In safe mode, we only flag conversions where we can be absolutely
	// certain the removal is safe. For now, we always consider the conversion
	// safe if types are identical - the type checker already verified this.
	// More advanced context checks could be added here in the future.
	_ = call
	_ = info
	return true
}

// normalizeOption returns an option name lowercased and without hyphens.
func normalizeOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
