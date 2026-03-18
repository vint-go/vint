package use_raw_string_regexp

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseRawStringRegexpRule suggests using raw string literals for regexp patterns
// that contain backslash escapes, since raw strings are easier to read.
type UseRawStringRegexpRule struct{}

// Apply applies the rule to given file.
func (r *UseRawStringRegexpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRawStringRegexp{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseRawStringRegexpRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintRawStringRegexp{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseRawStringRegexpRule) Name() string {
	return "useRawStringRegexp"
}

// Group returns the rule group.
func (*UseRawStringRegexpRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseRawStringRegexpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRawStringRegexp struct {
	onFailure func(lint.Failure)
}

func (w *lintRawStringRegexp) Visit(node ast.Node) ast.Visitor {
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

	funcName := sel.Sel.Name
	if funcName != "Compile" && funcName != "MustCompile" &&
		funcName != "CompilePOSIX" && funcName != "MustCompilePOSIX" &&
		funcName != "MatchString" && funcName != "Match" {
		return w
	}

	if len(call.Args) < 1 {
		return w
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	// Only flag interpreted strings (double-quoted), not raw strings (backtick-quoted)
	if !strings.HasPrefix(lit.Value, "\"") {
		return w
	}

	// Unquote to get the actual string value
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	// Check if the value contains backslashes (which means the interpreted string
	// had double backslashes like \\d that could be simplified to \d in a raw string)
	if !strings.Contains(value, `\`) {
		return w
	}

	// Cannot use raw string literal if the value contains backticks
	if strings.Contains(value, "`") {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    "regexp pattern can be simplified by using a raw string literal (backtick)",
	})

	return w
}
