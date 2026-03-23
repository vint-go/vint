package no_unkeyed_literal

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnkeyedLiteralRule flags composite literals of struct types that do not use field names.
type NoUnkeyedLiteralRule struct{}

// Apply applies the rule to given file.
func (r *NoUnkeyedLiteralRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoUnkeyedLiteral{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoUnkeyedLiteralRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnkeyedLiteral{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnkeyedLiteralRule) Name() string {
	return "noUnkeyedLiteral"
}

// Group returns the rule group.
func (*NoUnkeyedLiteralRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnkeyedLiteralRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoUnkeyedLiteralRule) RequiresTypecheck() bool {
	return true
}

type lintNoUnkeyedLiteral struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoUnkeyedLiteral) Visit(node ast.Node) ast.Visitor {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	// Skip composite literals with no elements
	if len(cl.Elts) == 0 {
		return w
	}

	// Check if any elements already use keys. If all use keys, it's fine.
	// If none use keys, it might be unkeyed.
	hasKeyedField := false
	for _, elt := range cl.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); ok {
			hasKeyedField = true
			break
		}
	}

	if hasKeyedField {
		// All or some fields are keyed; Go does not allow mixing, so this is keyed.
		return w
	}

	// Determine the underlying type of the composite literal
	t := w.pkg.TypeOf(cl)
	if t == nil {
		return w
	}

	// Unwrap pointer types
	underlying := t.Underlying()

	// Only flag struct types
	if _, ok := underlying.(*types.Struct); !ok {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       cl,
		Category:   lint.FailureCategoryStyle,
		Failure:    "composite literal uses unkeyed fields",
	})

	return w
}
