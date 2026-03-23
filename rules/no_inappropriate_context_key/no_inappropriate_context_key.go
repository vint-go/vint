package no_inappropriate_context_key

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInappropriateContextKeyRule detects calls to context.WithValue that use a
// built-in type (such as string or int) as the key, which can cause collisions
// between packages. Instead, an unexported custom type should be used.
type NoInappropriateContextKeyRule struct{}

// Apply applies the rule to given file.
func (r *NoInappropriateContextKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintInappropriateContextKey{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInappropriateContextKeyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInappropriateContextKey{
		file:      file,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInappropriateContextKeyRule) Name() string {
	return "noInappropriateContextKey"
}

// Group returns the rule group.
func (*NoInappropriateContextKeyRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInappropriateContextKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoInappropriateContextKeyRule) RequiresTypecheck() bool {
	return true
}

type lintInappropriateContextKey struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintInappropriateContextKey) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(call.Fun, "context", "WithValue") {
		return w
	}

	// context.WithValue takes 3 arguments: ctx, key, val
	if len(call.Args) != 3 {
		return w
	}

	info := w.file.Pkg.TypesInfo()
	if info == nil {
		return w
	}

	keyTypeAndValue, ok := info.Types[call.Args[1]]
	if !ok {
		return w
	}

	keyType := keyTypeAndValue.Type
	if keyType == nil {
		return w
	}

	if basic, ok := keyType.(*types.Basic); ok && basic.Kind() != types.Invalid {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("should not use basic type %s as key in context.WithValue", keyType),
		})
	}

	return w
}
