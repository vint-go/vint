package no_duplicate_option

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateOptionRule detects duplicated option function arguments in variadic function calls.
type NoDuplicateOptionRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateOptionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoDuplicateOption{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateOptionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoDuplicateOption{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateOptionRule) Name() string {
	return "noDuplicateOption"
}

// Group returns the rule group.
func (*NoDuplicateOptionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateOptionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoDuplicateOptionRule) RequiresTypecheck() bool {
	return true
}

type lintNoDuplicateOption struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoDuplicateOption) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	w.checkCall(call)
	return w
}

func (w *lintNoDuplicateOption) checkCall(call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	// Skip calls with ellipsis expansion (e.g., f(opts...))
	if call.Ellipsis != token.NoPos {
		return
	}

	funcType := w.pkg.TypeOf(call.Fun)
	if funcType == nil {
		return
	}

	sig, ok := funcType.Underlying().(*types.Signature)
	if !ok || !sig.Variadic() {
		return
	}

	last := sig.Params().Len() - 1
	sliceType, ok := sig.Params().At(last).Type().(*types.Slice)
	if !ok {
		return
	}

	if last > len(call.Args) {
		return
	}

	// Check if the variadic element type is a function type (option pattern)
	if !isOptionType(sliceType.Elem()) {
		return
	}

	// Check variadic args for duplicates
	variadicArgs := call.Args[last:]
	seen := make(map[string]bool)
	for _, arg := range variadicArgs {
		code := astutils.GoFmt(arg)
		if code == "" {
			continue
		}
		if seen[code] {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       arg,
				Failure:    "duplicated option function argument: " + code,
			})
		}
		seen[code] = true
	}
}

// isOptionType checks whether a type is a function type with at least one parameter (option pattern).
func isOptionType(t types.Type) bool {
	sig, ok := t.Underlying().(*types.Signature)
	if !ok {
		return false
	}
	return sig.Params().Len() > 0
}
