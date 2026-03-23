package no_invalid_flag_name

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidFlagNameRule detects suspicious flag names passed to Go's flag
// package functions. It warns when flag names are empty, start with a hyphen,
// contain equals signs, or contain whitespace characters.
type NoInvalidFlagNameRule struct{}

// flagFunctions lists the flag package functions that accept a flag name
// as their first argument.
var flagFunctions = []string{
	"Bool",
	"BoolVar",
	"Duration",
	"DurationVar",
	"Float64",
	"Float64Var",
	"Int",
	"IntVar",
	"Int64",
	"Int64Var",
	"String",
	"StringVar",
	"Uint",
	"UintVar",
	"Uint64",
	"Uint64Var",
}

// Apply applies the rule to given file.
func (r *NoInvalidFlagNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidFlagName{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidFlagNameRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInvalidFlagName{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidFlagNameRule) Name() string {
	return "noInvalidFlagName"
}

// Group returns the rule group.
func (*NoInvalidFlagNameRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidFlagNameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInvalidFlagName struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidFlagName) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fn := range flagFunctions {
		if !astutils.IsPkgDotName(callExpr.Fun, "flag", fn) {
			continue
		}

		// Determine which argument is the flag name.
		// For XxxVar functions the flag name is the second argument (index 1),
		// for others it is the first argument (index 0).
		nameArgIdx := 0
		if strings.HasSuffix(fn, "Var") {
			nameArgIdx = 1
		}

		if len(callExpr.Args) <= nameArgIdx {
			return w
		}

		lit, ok := callExpr.Args[nameArgIdx].(*ast.BasicLit)
		if !ok {
			return w
		}

		// Strip quotes from the string literal value.
		flagName := strings.Trim(lit.Value, "\"`")

		if msg := checkFlagName(flagName); msg != "" {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       callExpr,
				Category:   lint.FailureCategoryLogic,
				Failure:    msg,
			})
		}

		return w
	}

	return w
}

// checkFlagName returns an error message if the flag name is invalid,
// or empty string if it is valid.
func checkFlagName(name string) string {
	if name == "" {
		return "flag name is empty"
	}
	if strings.HasPrefix(name, "-") {
		return "flag name \"" + name + "\" should not start with a hyphen"
	}
	if strings.Contains(name, "=") {
		return "flag name \"" + name + "\" should not contain '='"
	}
	for _, r := range name {
		if unicode.IsSpace(r) {
			return "flag name \"" + name + "\" should not contain whitespace"
		}
	}
	return ""
}
