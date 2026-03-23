package no_dynamic_format_string

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDynamicFormatStringRule detects suspicious formatting calls where a dynamic
// (non-literal) string is passed as the format argument without additional arguments.
type NoDynamicFormatStringRule struct{}

// Apply applies the rule to given file.
func (r *NoDynamicFormatStringRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDynamicFormat{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDynamicFormatStringRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDynamicFormat{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDynamicFormatStringRule) Name() string {
	return "noDynamicFormatString"
}

// Group returns the rule group.
func (*NoDynamicFormatStringRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDynamicFormatStringRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDynamicFormat struct {
	onFailure func(lint.Failure)
}

// formatFuncs maps fmt function names to the index of the format string argument.
// For most functions it is 0 (first arg), for Fprintf/Fscanf it is 1 (second arg,
// after the io.Writer).
var formatFuncs = map[string]int{
	"Sprintf":  0,
	"Printf":   0,
	"Errorf":   0,
	"Fprintf":  1,
	"Fatalf":   0,
	"Panicf":   0,
	"Logf":     0,
	"Warnf":    0,
	"Infof":    0,
	"Debugf":   0,
	"Tracef":   0,
}

func (w *lintDynamicFormat) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for funcName, formatArgIdx := range formatFuncs {
		if !astutils.IsPkgDotName(ce.Fun, "fmt", funcName) {
			continue
		}

		// Check that the call has exactly formatArgIdx+1 arguments
		// (i.e., no additional format arguments beyond the format string itself).
		if len(ce.Args) != formatArgIdx+1 {
			return w
		}

		// Check if the format argument is a non-literal (dynamic) expression.
		formatArg := ce.Args[formatArgIdx]
		if _, isLit := formatArg.(*ast.BasicLit); isLit {
			return w // literal string is fine
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "dynamic string used as format argument without additional args; use fmt." + funcName[:len(funcName)-1] + " or add format arguments",
		})

		return w
	}

	return w
}
