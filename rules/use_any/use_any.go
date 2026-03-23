package use_any

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseAnyRule proposes to replace `interface{}` with its alias `any`.
type UseAnyRule struct{}

// Apply applies the rule to given file.
func (*UseAnyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := lintUseAny{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}
	fileAst := file.AST
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*UseAnyRule) Name() string {
	return "useAny"
}

// Group returns the rule group.
func (*UseAnyRule) Group() string {
	return "style"
}

type lintUseAny struct {
	onFailure func(lint.Failure)
}

func (w lintUseAny) Visit(n ast.Node) ast.Visitor {
	it, ok := n.(*ast.InterfaceType)
	if !ok {
		return w
	}

	if len(it.Methods.List) != 0 {
		return w // it is not and empty interface
	}

	w.onFailure(lint.Failure{
		Node:       n,
		Confidence: 1,
		Category:   lint.FailureCategoryNaming,
		Failure:    "since Go 1.18 'interface{}' can be replaced by 'any'",
	})

	return w
}

// CacheTier returns the cache tier for this rule.
func (*UseAnyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
