package no_redundant_type_assertion

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantTypeAssertionRule detects type assertions where the source and
// destination types are identical. When a value is asserted to a type it
// already possesses, the assertion is unnecessary and can be removed.
type NoRedundantTypeAssertionRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantTypeAssertionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantTypeAssertion{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantTypeAssertionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantTypeAssertion{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantTypeAssertionRule) Name() string {
	return "noRedundantTypeAssertion"
}

// Group returns the rule group.
func (*NoRedundantTypeAssertionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantTypeAssertionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*NoRedundantTypeAssertionRule) RequiresTypecheck() bool {
	return true
}

type lintRedundantTypeAssertion struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintRedundantTypeAssertion) Visit(node ast.Node) ast.Visitor {
	ta, ok := node.(*ast.TypeAssertExpr)
	if !ok {
		return w
	}

	// Type switch assertions (ta.Type == nil means x.(type)) are not checked.
	if ta.Type == nil {
		return w
	}

	typesInfo := w.pkg.TypesInfo()
	if typesInfo == nil {
		return w
	}

	// Get the type of the expression being asserted (the LHS of the dot).
	exprType := w.pkg.TypeOf(ta.X)
	if exprType == nil {
		return w
	}

	// Get the asserted-to type (the type inside the parentheses).
	assertedType := w.pkg.TypeOf(ta.Type)
	if assertedType == nil {
		return w
	}

	// If the types are identical, the assertion is redundant.
	if !types.Identical(exprType, assertedType) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       ta,
		Failure:    fmt.Sprintf("redundant type assertion: expression already has type %s", assertedType.String()),
	})

	return w
}
