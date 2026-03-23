package no_useless_math_call

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUselessMathCallRule detects calls to math.Ceil, math.Floor, and math.Round
// on floats converted from integers, which have no effect.
type NoUselessMathCallRule struct{}

// mathFuncs lists the math functions that are no-ops on integer-converted values.
var mathFuncs = []string{"Ceil", "Floor", "Round", "Trunc", "RoundToEven"}

// Apply applies the rule to given file.
func (r *NoUselessMathCallRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintUselessMathCall{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUselessMathCallRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintUselessMathCall{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoUselessMathCallRule) Name() string {
	return "noUselessMathCall"
}

// Group returns the rule group.
func (*NoUselessMathCallRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUselessMathCallRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoUselessMathCallRule) RequiresTypecheck() bool {
	return true
}

type lintUselessMathCall struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintUselessMathCall) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if it's a call to one of the targeted math functions
	var funcName string
	for _, name := range mathFuncs {
		if astutils.IsPkgDotName(call.Fun, "math", name) {
			funcName = name
			break
		}
	}

	if funcName == "" {
		return w
	}

	// Must have exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	arg := call.Args[0]

	// Check if the argument is a float64() conversion of an integer
	innerCall, ok := arg.(*ast.CallExpr)
	if !ok {
		return w
	}

	// The conversion target must be float64
	if !astutils.IsIdent(innerCall.Fun, "float64") {
		return w
	}

	if len(innerCall.Args) != 1 {
		return w
	}

	// Check if the inner argument has an integer type
	innerArg := innerCall.Args[0]
	innerType := w.pkg.TypeOf(innerArg)
	if innerType == nil {
		return w
	}

	if !isIntegerType(innerType) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       call,
		Failure:    fmt.Sprintf("math.%s called on a float64 converted from an integer, which is pointless", funcName),
	})

	return w
}

// isIntegerType checks if the type is an integer type (signed or unsigned).
func isIntegerType(t types.Type) bool {
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Info()&types.IsInteger != 0
}
