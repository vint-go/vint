package no_invalid_errors_as

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInvalidErrorsAsRule checks that the second argument to errors.As is a
// pointer to a type implementing the error interface or a pointer to any
// interface type.
type NoInvalidErrorsAsRule struct{}

// Apply applies the rule to given file.
func (*NoInvalidErrorsAsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoInvalidErrorsAs{file: file, onFailure: onFailure}
	if w.file.Pkg.TypeCheck() != nil {
		return nil
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoInvalidErrorsAsRule) Name() string {
	return "noInvalidErrorsAs"
}

// Group returns the rule group.
func (*NoInvalidErrorsAsRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidErrorsAsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoInvalidErrorsAsRule) RequiresTypecheck() bool {
	return true
}

type lintNoInvalidErrorsAs struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoInvalidErrorsAs) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(call.Fun, "errors", "As") {
		return w
	}

	if len(call.Args) < 2 {
		return w
	}

	secondArg := call.Args[1]
	argType := w.file.Pkg.TypeOf(secondArg)
	if argType == nil {
		return w
	}

	if !isValidErrorsAsTarget(argType) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryErrors,
			Confidence: 1,
			Node:       secondArg,
			Failure:    fmt.Sprintf("second argument to errors.As must be a pointer to an interface or a type implementing error, got %s", argType),
		})
	}

	return w
}

// isValidErrorsAsTarget checks if the type is a valid target for errors.As.
// The target must be a non-nil pointer to either:
//   - any interface type, or
//   - a type that implements the error interface.
//
// Crucially, if Error() is defined only on a pointer receiver (*T), then
// the target must be **T (a pointer to a *T variable), not *T.
func isValidErrorsAsTarget(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}

	elem := ptr.Elem()

	// If the element is an interface type, it's always valid.
	if _, ok := elem.Underlying().(*types.Interface); ok {
		return true
	}

	// The element type itself must implement error.
	// Note: we do NOT check *elem here. If only *elem implements error
	// (pointer receiver), then the target should be **elem, not *elem.
	errorIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(elem, errorIface)
}
