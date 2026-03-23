package use_infinite_for

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseInfiniteForRule detects `for true { ... }` and suggests using `for { ... }` instead.
type UseInfiniteForRule struct{}

// Apply applies the rule to given file.
func (r *UseInfiniteForRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseInfiniteFor{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseInfiniteForRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseInfiniteFor{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseInfiniteForRule) Name() string {
	return "useInfiniteFor"
}

// Group returns the rule group.
func (*UseInfiniteForRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseInfiniteForRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseInfiniteFor struct {
	onFailure func(lint.Failure)
}

func (w *lintUseInfiniteFor) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// We only flag `for true { ... }` — a for statement whose condition
	// is the boolean literal `true` and that has no init or post statements.
	if forStmt.Init != nil || forStmt.Post != nil {
		return w
	}

	ident, ok := forStmt.Cond.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != "true" {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       forStmt,
		Failure:    "use for { ... } instead of for true { ... }",
	})

	return w
}
