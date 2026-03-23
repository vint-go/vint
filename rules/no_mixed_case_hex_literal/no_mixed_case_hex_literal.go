package no_mixed_case_hex_literal

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMixedCaseHexLiteralRule detects hex literals that have mixed case letter digits
// or use the uppercase 0X prefix.
type NoMixedCaseHexLiteralRule struct{}

// Apply applies the rule to given file.
func (r *NoMixedCaseHexLiteralRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoMixedCaseHex{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMixedCaseHexLiteralRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintNoMixedCaseHex{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMixedCaseHexLiteralRule) Name() string {
	return "noMixedCaseHexLiteral"
}

// Group returns the rule group.
func (*NoMixedCaseHexLiteralRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMixedCaseHexLiteralRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoMixedCaseHex struct {
	onFailure func(lint.Failure)
}

func (w *lintNoMixedCaseHex) Visit(node ast.Node) ast.Visitor {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind != token.INT {
		return w
	}

	value := lit.Value

	// Check if it's a hex literal (starts with 0x or 0X)
	if len(value) < 3 {
		return w
	}
	if value[0] != '0' || (value[1] != 'x' && value[1] != 'X') {
		return w
	}

	// Check for uppercase prefix 0X
	if value[1] == 'X' {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       lit,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("hex literal %s uses uppercase 0X prefix, use lowercase 0x instead", value),
		})
		return w
	}

	// Check for mixed case in hex digit letters (after 0x prefix)
	digits := value[2:]
	// Remove underscores (Go allows _ in numeric literals)
	digits = strings.ReplaceAll(digits, "_", "")

	hasUpper := false
	hasLower := false
	for _, ch := range digits {
		if ch >= 'a' && ch <= 'f' {
			hasLower = true
		}
		if ch >= 'A' && ch <= 'F' {
			hasUpper = true
		}
	}

	if hasUpper && hasLower {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       lit,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("hex literal %s has mixed case letter digits, use consistent casing", value),
		})
	}

	return w
}
