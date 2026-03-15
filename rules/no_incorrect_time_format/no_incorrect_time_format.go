package no_incorrect_time_format

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIncorrectTimeFormatRule checks for the use of incorrect time format strings
// in calls to time.Format and time.Parse.
type NoIncorrectTimeFormatRule struct{}

// Apply applies the rule to given file.
func (r *NoIncorrectTimeFormatRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintIncorrectTimeFormat{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoIncorrectTimeFormatRule) Name() string {
	return "noIncorrectTimeFormat"
}

// Group returns the rule group.
func (*NoIncorrectTimeFormatRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoIncorrectTimeFormatRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIncorrectTimeFormat struct {
	onFailure func(lint.Failure)
}

func (w *lintIncorrectTimeFormat) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	var layoutArg ast.Expr

	switch fun := ce.Fun.(type) {
	case *ast.SelectorExpr:
		// Check for time.Parse(layout, value) — package-level function call
		if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "time" && fun.Sel.Name == "Parse" {
			if len(ce.Args) >= 1 {
				layoutArg = ce.Args[0]
			}
		}
		// Check for t.Format(layout) — method call on any receiver with method name "Format"
		if fun.Sel.Name == "Format" {
			if len(ce.Args) >= 1 {
				layoutArg = ce.Args[0]
			}
		}
	}

	if layoutArg == nil {
		return w
	}

	// Only check string literals
	lit, ok := layoutArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	// Remove quotes from the string literal value
	value := lit.Value
	if len(value) >= 2 {
		value = value[1 : len(value)-1]
	}

	if msg := checkTimeFormat(value); msg != "" {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    msg,
		})
	}

	return w
}

// checkTimeFormat checks a time format string for common mistakes and returns
// a diagnostic message if a problem is found, or empty string if the format is ok.
func checkTimeFormat(format string) string {
	// Check for swapped month and day: 2006-02-01 instead of 2006-01-02
	if strings.Contains(format, "2006-02-01") {
		return "month and day are swapped in time format string: use \"2006-01-02\" not \"2006-02-01\""
	}

	// Check for swapped month and day with slash separator: 2006/02/01 instead of 2006/01/02
	if strings.Contains(format, "2006/02/01") {
		return "month and day are swapped in time format string: use \"2006/01/02\" not \"2006/02/01\""
	}

	return ""
}
