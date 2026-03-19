package no_context_keys_type

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// ContextKeysType disallows the usage of basic types in [context.WithValue].
type ContextKeysType struct{}

// Apply applies the rule to given file.
func (*ContextKeysType) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	walker := lintContextKeyTypes{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck()
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*ContextKeysType) Name() string {
	return "noContextKeysType"
}

// Group returns the rule group.
func (*ContextKeysType) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*ContextKeysType) RequiresTypecheck() bool { return true }

type lintContextKeyTypes struct {
	file      *lint.File
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

func (w lintContextKeyTypes) Visit(n ast.Node) ast.Visitor {
	if n, ok := n.(*ast.CallExpr); ok {
		checkContextKeyType(w, n)
	}

	return w
}

func checkContextKeyType(w lintContextKeyTypes, x *ast.CallExpr) {
	f := w.file
	if !astutils.IsPkgDotName(x.Fun, "context", "WithValue") {
		return
	}

	// key is second argument to context.WithValue
	if len(x.Args) != 3 {
		return
	}
	key := f.Pkg.TypesInfo().Types[x.Args[1]]

	if ktyp, ok := key.Type.(*types.Basic); ok && ktyp.Kind() != types.Invalid {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       x,
			Category:   lint.FailureCategoryContent,
			Failure:    fmt.Sprintf("should not use basic type %s as key in context.WithValue", key.Type),
		})
	}
}

// CacheTier returns the cache tier for this rule.
func (*ContextKeysType) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
