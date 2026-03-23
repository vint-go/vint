package no_invalid_regexp

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp/syntax"
	"strconv"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidRegexpRule validates that arguments to regexp.Compile, regexp.MustCompile,
// regexp.Match, and related functions are valid regular expressions.
type NoInvalidRegexpRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidRegexpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidRegexp{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidRegexpRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInvalidRegexp{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidRegexpRule) Name() string {
	return "noInvalidRegexp"
}

// Group returns the rule group.
func (*NoInvalidRegexpRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidRegexpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// regexpFunctions lists the regexp package functions that accept a pattern string
// as their first argument, along with the syntax flags to use for validation.
var regexpFunctions = map[string]syntax.Flags{
	"Compile":      syntax.Perl,
	"MustCompile":  syntax.Perl,
	"Match":        syntax.Perl,
	"MatchString":  syntax.Perl,
	"MatchReader":  syntax.Perl,
	"CompilePOSIX": syntax.POSIX,
}

type lintInvalidRegexp struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidRegexp) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	if pkgIdent.Name != "regexp" {
		return w
	}

	flags, ok := regexpFunctions[sel.Sel.Name]
	if !ok {
		return w
	}

	if len(call.Args) < 1 {
		return w
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	pattern, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	if _, err := syntax.Parse(pattern, flags); err != nil {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("invalid regular expression: %s", err),
		})
	}

	return w
}
