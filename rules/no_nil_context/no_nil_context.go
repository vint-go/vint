package no_nil_context

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNilContextRule detects calls where nil is passed as a context.Context argument.
type NoNilContextRule struct{}

// Apply applies the rule to given file.
func (r *NoNilContextRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNilContext{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNilContextRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNilContext{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNilContextRule) Name() string {
	return "noNilContext"
}

// Group returns the rule group.
func (*NoNilContextRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNilContextRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoNilContextRule) RequiresTypecheck() bool {
	return true
}

type lintNilContext struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNilContext) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	w.checkCall(call)
	return w
}

func (w *lintNilContext) checkCall(call *ast.CallExpr) {
	info := w.pkg.TypesInfo()
	if info == nil {
		return
	}

	// Get the type of the function being called.
	funcType := w.pkg.TypeOf(call.Fun)
	if funcType == nil {
		return
	}

	sig, ok := funcType.Underlying().(*types.Signature)
	if !ok {
		return
	}

	params := sig.Params()
	if params == nil {
		return
	}

	// Check each argument against the corresponding parameter.
	for i, arg := range call.Args {
		// Skip variadic expansion arguments beyond parameter count.
		paramIdx := i
		if paramIdx >= params.Len() {
			if sig.Variadic() {
				paramIdx = params.Len() - 1
			} else {
				break
			}
		}

		// Check if the argument is nil.
		ident, ok := arg.(*ast.Ident)
		if !ok || ident.Name != "nil" {
			continue
		}

		// Check if the parameter type is context.Context.
		paramType := params.At(paramIdx).Type()
		if isContextType(paramType) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       arg,
				Failure:    "nil context passed, use context.TODO or context.Background instead",
			})
		}
	}
}

// isContextType checks whether a type is context.Context.
func isContextType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	return obj.Pkg().Path() == "context" && obj.Name() == "Context"
}
