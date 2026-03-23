package no_multi_line_func_break

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMultiLineFuncBreakRule detects multi-line function signatures where the
// parameter list does not start on the same line as the func keyword.
type NoMultiLineFuncBreakRule struct {
	multiFunc bool
}

// Configure validates and applies the rule configuration.
func (r *NoMultiLineFuncBreakRule) Configure(arguments lint.Arguments) error {
	r.multiFunc = false

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMultiLineFuncBreak" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == normalizeOption("multi-func") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for multi-func in "noMultiLineFuncBreak" rule; need bool but got %T`, v)
			}
			r.multiFunc = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoMultiLineFuncBreakRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !r.multiFunc {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintMultiLineFuncBreak{
		onFailure: onFailure,
		file:      file,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMultiLineFuncBreakRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if !r.multiFunc {
		return nil
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintMultiLineFuncBreak{onFailure: onFailure, file: file}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMultiLineFuncBreakRule) Name() string {
	return "noMultiLineFuncBreak"
}

// Group returns the rule group.
func (*NoMultiLineFuncBreakRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMultiLineFuncBreakRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMultiLineFuncBreak struct {
	onFailure func(lint.Failure)
	file      *lint.File
}

func (w *lintMultiLineFuncBreak) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	if funcDecl.Type == nil || funcDecl.Type.Params == nil {
		return w
	}

	params := funcDecl.Type.Params
	if len(params.List) == 0 {
		return w
	}

	// Get the line of the func keyword (the start of the FuncDecl).
	funcLine := w.file.ToPosition(funcDecl.Pos()).Line
	// Get the line of the first parameter.
	firstParamLine := w.file.ToPosition(params.List[0].Pos()).Line

	if firstParamLine > funcLine {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       funcDecl,
			Category:   lint.FailureCategoryStyle,
			Failure:    "multi-line function signature should have the first parameter on the same line as the func keyword",
		})
	}

	return w
}

// normalizeOption returns an option name lowercased and without hyphens.
func normalizeOption(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}
