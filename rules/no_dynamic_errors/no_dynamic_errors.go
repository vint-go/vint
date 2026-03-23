package no_dynamic_errors

import (
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/lint"
)

// NoDynamicErrorsRule flags the creation of dynamic errors inside functions.
type NoDynamicErrorsRule struct{}

// Apply applies the rule to given file.
func (*NoDynamicErrorsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDynamicErrors{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoDynamicErrorsRule) Name() string {
	return "noDynamicErrors"
}

// Group returns the rule group.
func (*NoDynamicErrorsRule) Group() string {
	return "correctness"
}

type lintNoDynamicErrors struct {
	onFailure func(lint.Failure)
}

func (w *lintNoDynamicErrors) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "errors", "New") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryErrors,
			Failure:    "use of errors.New() inside a function: define sentinel errors at package level",
		})
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "fmt", "Errorf") {
		if !hasWrappingVerb(ce) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryErrors,
				Failure:    "use of fmt.Errorf() without %w: use fmt.Errorf with %w to wrap sentinel errors",
			})
		}
		return w
	}

	return w
}

// hasWrappingVerb checks if a fmt.Errorf call contains %w in its format string.
func hasWrappingVerb(ce *ast.CallExpr) bool {
	if len(ce.Args) == 0 {
		return false
	}

	formatArg, ok := ce.Args[0].(*ast.BasicLit)
	if !ok {
		return false
	}

	return strings.Contains(formatArg.Value, "%w")
}
