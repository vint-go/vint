package no_redundant_switch_true

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantSwitchTrueRule detects switch statements with an explicit `true` tag
// that is redundant because `switch { ... }` is equivalent to `switch true { ... }`.
type NoRedundantSwitchTrueRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantSwitchTrueRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoRedundantSwitchTrue{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantSwitchTrueRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoRedundantSwitchTrue{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoRedundantSwitchTrueRule) Name() string {
	return "noRedundantSwitchTrue"
}

// Group returns the rule group.
func (*NoRedundantSwitchTrueRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantSwitchTrueRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoRedundantSwitchTrue struct {
	onFailure func(lint.Failure)
}

func (w *lintNoRedundantSwitchTrue) Visit(node ast.Node) ast.Visitor {
	switchStmt, ok := node.(*ast.SwitchStmt)
	if !ok {
		return w
	}

	// Check if the switch tag is the identifier `true`
	ident, ok := switchStmt.Tag.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name == "true" {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    "redundant switch true; use switch without a tag instead",
			Node:       switchStmt,
		})
	}

	return w
}
