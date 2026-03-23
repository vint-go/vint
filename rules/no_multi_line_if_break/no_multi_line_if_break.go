package no_multi_line_if_break

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMultiLineIfBreakRule detects multi-line if conditions where the condition
// does not start on the same line as the if keyword.
type NoMultiLineIfBreakRule struct {
	multiIf bool
}

// Configure validates and applies the rule configuration.
func (r *NoMultiLineIfBreakRule) Configure(arguments lint.Arguments) error {
	r.multiIf = false

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMultiLineIfBreak" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == normalizeOption("multi-if") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for multi-if in "noMultiLineIfBreak" rule; need bool but got %T`, v)
			}
			r.multiIf = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoMultiLineIfBreakRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !r.multiIf {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintMultiLineIfBreak{
		onFailure: onFailure,
		file:      file,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMultiLineIfBreakRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if !r.multiIf {
		return nil
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintMultiLineIfBreak{onFailure: onFailure, file: file}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMultiLineIfBreakRule) Name() string {
	return "noMultiLineIfBreak"
}

// Group returns the rule group.
func (*NoMultiLineIfBreakRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMultiLineIfBreakRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMultiLineIfBreak struct {
	onFailure func(lint.Failure)
	file      *lint.File
}

func (w *lintMultiLineIfBreak) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	if ifStmt.Cond == nil {
		return w
	}

	ifLine := w.file.ToPosition(ifStmt.Pos()).Line
	condLine := w.file.ToPosition(ifStmt.Cond.Pos()).Line

	if condLine > ifLine {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ifStmt,
			Category:   lint.FailureCategoryStyle,
			Failure:    "multi-line if condition should start on the same line as the if keyword",
		})
	}

	return w
}

// normalizeOption returns an option name lowercased and without hyphens.
func normalizeOption(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}
