package use_error_method

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseErrorMethodRule detects fmt.Sprintf("%s", err) calls where err implements
// the error interface, and suggests using err.Error() instead.
type UseErrorMethodRule struct{}

// Apply applies the rule to given file.
func (r *UseErrorMethodRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintErrorMethod{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseErrorMethodRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintErrorMethod{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseErrorMethodRule) Name() string {
	return "useErrorMethod"
}

// Group returns the rule group.
func (*UseErrorMethodRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseErrorMethodRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*UseErrorMethodRule) RequiresTypecheck() bool {
	return true
}

type lintErrorMethod struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintErrorMethod) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for fmt.Sprintf("%s", errArg)
	if !astutils.IsPkgDotName(call.Fun, "fmt", "Sprintf") {
		return w
	}

	// Must have exactly 2 arguments: format string and one value
	if len(call.Args) != 2 {
		return w
	}

	// First argument must be the format string "%s"
	formatLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return w
	}

	if formatLit.Value != `"%s"` {
		return w
	}

	// Second argument must implement the error interface
	arg := call.Args[1]
	argType := w.pkg.TypeOf(arg)
	if argType == nil {
		return w
	}

	if !isErrorType(argType) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    `use err.Error() instead of fmt.Sprintf("%s", err)`,
	})

	return w
}

// isErrorType checks if the type implements the error interface.
func isErrorType(t types.Type) bool {
	// Check if the type is exactly the error interface or implements it
	errorInterface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(t, errorInterface) || types.Implements(types.NewPointer(t), errorInterface)
}
