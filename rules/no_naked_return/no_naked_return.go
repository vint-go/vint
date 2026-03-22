package no_naked_return

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNakedReturnRule flags naked return statements in functions that exceed a
// configurable line length threshold.
type NoNakedReturnRule struct {
	maxFuncLines int
}

const defaultMaxFuncLines = 30

// Configure validates and applies the rule configuration.
func (r *NoNakedReturnRule) Configure(arguments lint.Arguments) error {
	r.maxFuncLines = defaultMaxFuncLines

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int argument
		lines, ok := toInt(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "noNakedReturn" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.maxFuncLines = lines
		return nil
	}

	for k, v := range argKV {
		if normalizeOption(k) == normalizeOption("maxFuncLines") {
			lines, ok := toInt(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for maxFuncLines in "noNakedReturn" rule; need integer but got %T`, v)
			}
			r.maxFuncLines = lines
		}
	}

	return nil
}

// Apply applies the rule to the given file.
func (r *NoNakedReturnRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.maxFuncLines <= 0 {
		r.maxFuncLines = defaultMaxFuncLines
	}

	var failures []lint.Failure

	w := &lintNakedReturn{
		maxFuncLines: r.maxFuncLines,
		file:         file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNakedReturnRule) Name() string {
	return "noNakedReturn"
}

// Group returns the rule group.
func (*NoNakedReturnRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoNakedReturnRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNakedReturn struct {
	maxFuncLines int
	file         *lint.File
	onFailure    func(lint.Failure)
}

func (w *lintNakedReturn) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFunc(n.Type, n.Body)
	case *ast.FuncLit:
		w.checkFunc(n.Type, n.Body)
	}
	return w
}

func (w *lintNakedReturn) checkFunc(funcType *ast.FuncType, body *ast.BlockStmt) {
	if body == nil || funcType.Results == nil {
		return
	}

	// Check if function has named return parameters.
	hasNamedResults := false
	for _, field := range funcType.Results.List {
		if len(field.Names) > 0 {
			hasNamedResults = true
			break
		}
	}
	if !hasNamedResults {
		return
	}

	// Calculate the function body line count.
	startLine := w.file.ToPosition(body.Lbrace).Line
	endLine := w.file.ToPosition(body.Rbrace).Line
	bodyLines := endLine - startLine + 1

	if bodyLines <= w.maxFuncLines {
		return // function is short enough, naked returns are acceptable
	}

	// Collect named return parameter names for the failure message.
	var names []string
	for _, field := range funcType.Results.List {
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}

	// Find naked returns within this function body (but not in nested func literals).
	finder := &nakedReturnFinder{
		names:     names,
		onFailure: w.onFailure,
	}
	ast.Walk(finder, body)
}

type nakedReturnFinder struct {
	names     []string
	onFailure func(lint.Failure)
}

func (f *nakedReturnFinder) Visit(node ast.Node) ast.Visitor {
	// Skip nested function literals; they will be handled by the outer walker.
	if _, ok := node.(*ast.FuncLit); ok {
		return nil
	}

	rs, ok := node.(*ast.ReturnStmt)
	if !ok {
		return f
	}

	if len(rs.Results) > 0 {
		return f
	}

	f.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryStyle,
		Node:       rs,
		Failure:    fmt.Sprintf("naked return in func with %s as named returns, use explicit return values", strings.Join(f.names, ", ")),
	})

	return f
}

func normalizeOption(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}

// toInt converts a value to int, supporting both int and int64 (YAML v3 decodes
// integers as Go int, while some code paths may provide int64).
func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}
