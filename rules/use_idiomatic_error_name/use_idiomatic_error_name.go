package use_idiomatic_error_name

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseIdiomaticErrorNameRule lints naming of error variables.
// Error variables should be named `err` or have a prefix of `err` or `Err`
// (for exported variables). Sentinel errors (package-level error variables)
// should be named `ErrFoo`, not `ErrorFoo` or `FooError`.
type UseIdiomaticErrorNameRule struct{}

// Apply applies the rule to given file.
func (r *UseIdiomaticErrorNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintErrorName{
		fileAst:   file.AST,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseIdiomaticErrorNameRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintErrorName{
		fileAst:   file.AST,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseIdiomaticErrorNameRule) Name() string {
	return "useIdiomaticErrorName"
}

// Group returns the rule group.
func (*UseIdiomaticErrorNameRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseIdiomaticErrorNameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintErrorName struct {
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

func (w *lintErrorName) Visit(_ ast.Node) ast.Visitor {
	for _, decl := range w.fileAst.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			spec := spec.(*ast.ValueSpec)
			if len(spec.Names) != 1 || len(spec.Values) != 1 {
				continue
			}
			ce, ok := spec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			if !astutils.IsPkgDotName(ce.Fun, "errors", "New") && !astutils.IsPkgDotName(ce.Fun, "fmt", "Errorf") {
				continue
			}

			id := spec.Names[0]
			if id.Name == "_" {
				continue
			}

			prefix := "err"
			if id.IsExported() {
				prefix = "Err"
			}
			if !strings.HasPrefix(id.Name, prefix) {
				w.onFailure(lint.Failure{
					Node:       id,
					Confidence: 0.9,
					Category:   lint.FailureCategoryNaming,
					Failure:    fmt.Sprintf("error var %s should have name of the form %sFoo", id.Name, prefix),
				})
			}
		}
	}
	return nil
}
