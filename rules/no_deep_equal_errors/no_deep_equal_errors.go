package no_deep_equal_errors

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDeepEqualErrorsRule checks for the use of reflect.DeepEqual with error values.
type NoDeepEqualErrorsRule struct{}

// Apply applies the rule to given file.
func (r *NoDeepEqualErrorsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintNoDeepEqualErrors{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoDeepEqualErrorsRule) Name() string {
	return "noDeepEqualErrors"
}

// Group returns the rule group.
func (*NoDeepEqualErrorsRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeepEqualErrorsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoDeepEqualErrorsRule) RequiresTypecheck() bool {
	return true
}

type lintNoDeepEqualErrors struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoDeepEqualErrors) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "reflect", "DeepEqual") {
		return w
	}

	if len(ce.Args) != 2 {
		return w
	}

	arg0Type := w.file.Pkg.TypeOf(ce.Args[0])
	arg1Type := w.file.Pkg.TypeOf(ce.Args[1])

	if arg0Type != nil && implementsError(arg0Type) || arg1Type != nil && implementsError(arg1Type) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       ce,
			Failure:    "avoid using reflect.DeepEqual with error values, use errors.Is instead",
		})
	}

	return w
}

// errorInterface is the error interface type.
var errorInterface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

// implementsError returns true if the given type implements the error interface.
func implementsError(t types.Type) bool {
	return types.Implements(t, errorInterface) || types.Implements(types.NewPointer(t), errorInterface)
}
