package no_redundant_sprint

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantSprintRule detects redundant fmt.Sprint calls on values that are already strings.
type NoRedundantSprintRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantSprintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantSprint{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantSprintRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantSprint{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantSprintRule) Name() string {
	return "noRedundantSprint"
}

// Group returns the rule group.
func (*NoRedundantSprintRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantSprintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoRedundantSprintRule) RequiresTypecheck() bool {
	return true
}

type lintRedundantSprint struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintRedundantSprint) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for fmt.Sprint(stringArg) - single arg that is already a string
	if astutils.IsPkgDotName(call.Fun, "fmt", "Sprint") {
		w.checkSprintSingleStringArg(call)
		return w
	}

	// Check for fmt.Sprintf("%s", stringArg) - single %s with a string arg
	if astutils.IsPkgDotName(call.Fun, "fmt", "Sprintf") {
		w.checkSprintfRedundant(call)
		return w
	}

	return w
}

// checkSprintSingleStringArg checks if fmt.Sprint is called with a single string argument.
func (w *lintRedundantSprint) checkSprintSingleStringArg(call *ast.CallExpr) {
	if len(call.Args) != 1 {
		return
	}

	arg := call.Args[0]
	argType := w.pkg.TypeOf(arg)
	if argType == nil {
		return
	}

	if !isStringType(argType) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    "redundant fmt.Sprint call on a value that is already a string",
	})
}

// checkSprintfRedundant checks if fmt.Sprintf is called with just "%s" and a single string argument.
func (w *lintRedundantSprint) checkSprintfRedundant(call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}

	// First argument must be the format string "%s"
	formatLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	// Check if the format string is exactly "%s" (including quotes)
	if formatLit.Value != `"%s"` {
		return
	}

	// Second argument must be a string type
	arg := call.Args[1]
	argType := w.pkg.TypeOf(arg)
	if argType == nil {
		return
	}

	if !isStringType(argType) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    "redundant fmt.Sprintf call with %s on a value that is already a string",
	})
}

// isStringType checks if the type is the builtin string type.
func isStringType(t types.Type) bool {
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.String
}
