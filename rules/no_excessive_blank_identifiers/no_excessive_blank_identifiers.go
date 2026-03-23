package no_excessive_blank_identifiers

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExcessiveBlankIdentifiersRule checks assignments with too many blank identifiers.
type NoExcessiveBlankIdentifiersRule struct {
	maxBlankIdentifiers int
}

const defaultMaxBlankIdentifiers = 2

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoExcessiveBlankIdentifiersRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.maxBlankIdentifiers = defaultMaxBlankIdentifiers
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		val, ok := lint.ToInt64(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "noExcessiveBlankIdentifiers" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.maxBlankIdentifiers = int(val)
		return nil
	}

	r.maxBlankIdentifiers = defaultMaxBlankIdentifiers
	for k, v := range argKV {
		if isRuleOption(k, "max-blank-identifiers") {
			val, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for max-blank-identifiers in "noExcessiveBlankIdentifiers" rule; need integer but got %T`, v)
			}
			r.maxBlankIdentifiers = int(val)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoExcessiveBlankIdentifiersRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	maxBlanks := r.maxBlankIdentifiers
	if maxBlanks == 0 {
		maxBlanks = defaultMaxBlankIdentifiers
	}

	var failures []lint.Failure

	w := &lintBlankIdentifiers{
		maxBlanks: maxBlanks,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoExcessiveBlankIdentifiersRule) Name() string {
	return "noExcessiveBlankIdentifiers"
}

// Group returns the rule group.
func (*NoExcessiveBlankIdentifiersRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoExcessiveBlankIdentifiersRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBlankIdentifiers struct {
	maxBlanks int
	onFailure func(lint.Failure)
}

func (w *lintBlankIdentifiers) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	blankCount := 0
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if ok && ident.Name == "_" {
			blankCount++
		}
	}

	if blankCount > w.maxBlanks {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("assignment has too many blank identifiers (%d > %d)", blankCount, w.maxBlanks),
			Node:       assign,
		})
	}

	return w
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
