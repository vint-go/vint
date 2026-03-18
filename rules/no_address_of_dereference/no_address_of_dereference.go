package no_address_of_dereference

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoAddressOfDereferenceRule detects taking the address of a dereferenced pointer (&*x).
type NoAddressOfDereferenceRule struct{}

// Apply applies the rule to given file.
func (r *NoAddressOfDereferenceRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoAddressOfDeref{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoAddressOfDereferenceRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoAddressOfDeref{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoAddressOfDereferenceRule) Name() string {
	return "noAddressOfDereference"
}

// Group returns the rule group.
func (*NoAddressOfDereferenceRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoAddressOfDereferenceRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoAddressOfDeref struct {
	onFailure func(lint.Failure)
}

func (w *lintNoAddressOfDeref) Visit(node ast.Node) ast.Visitor {
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		return w
	}

	// Check for &(*x) pattern: UnaryExpr with Op == & and operand is StarExpr
	if unary.Op != token.AND {
		return w
	}

	_, ok = unary.X.(*ast.StarExpr)
	if !ok {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       unary,
		Category:   lint.FailureCategoryLogic,
		Failure:    "&*x does not copy x, use a local variable to copy the value",
	})

	return w
}
