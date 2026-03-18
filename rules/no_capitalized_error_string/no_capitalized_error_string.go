package no_capitalized_error_string

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoCapitalizedErrorStringRule checks that error strings are not capitalized
// (unless beginning with proper nouns or acronyms) and do not end with punctuation.
type NoCapitalizedErrorStringRule struct{}

// Apply applies the rule to given file.
func (r *NoCapitalizedErrorStringRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoCapErrStr{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoCapitalizedErrorStringRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintNoCapErrStr{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoCapitalizedErrorStringRule) Name() string {
	return "noCapitalizedErrorString"
}

// Group returns the rule group.
func (*NoCapitalizedErrorStringRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoCapitalizedErrorStringRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoCapErrStr struct {
	onFailure func(lint.Failure)
}

func (w *lintNoCapErrStr) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if len(ce.Args) < 1 {
		return w
	}

	if !isErrorStringFunc(ce) {
		return w
	}

	str := getStringArg(ce)
	if str == nil {
		return w
	}

	s, err := strconv.Unquote(str.Value)
	if err != nil || s == "" {
		return w
	}

	if msg, ok := checkErrorString(s); !ok {
		w.onFailure(lint.Failure{
			Confidence: 0.8,
			Node:       str,
			Category:   lint.FailureCategoryStyle,
			Failure:    msg,
		})
	}

	return w
}

// isErrorStringFunc checks if the call is errors.New, fmt.Errorf, or errors.Wrap/Wrapf/WithMessage/WithMessagef/Errorf.
func isErrorStringFunc(ce *ast.CallExpr) bool {
	if astutils.IsPkgDotName(ce.Fun, "errors", "New") ||
		astutils.IsPkgDotName(ce.Fun, "fmt", "Errorf") ||
		astutils.IsPkgDotName(ce.Fun, "errors", "Errorf") ||
		astutils.IsPkgDotName(ce.Fun, "errors", "Wrap") ||
		astutils.IsPkgDotName(ce.Fun, "errors", "Wrapf") ||
		astutils.IsPkgDotName(ce.Fun, "errors", "WithMessage") ||
		astutils.IsPkgDotName(ce.Fun, "errors", "WithMessagef") {
		return true
	}
	return false
}

// getStringArg returns the first string literal argument from the call expression.
// For functions like Wrap/Wrapf/WithMessage/WithMessagef, the message is the second argument.
func getStringArg(ce *ast.CallExpr) *ast.BasicLit {
	// Try first argument
	if lit, ok := ce.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return lit
	}
	// For Wrap-style functions, try second argument
	if len(ce.Args) >= 2 {
		if lit, ok := ce.Args[1].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			return lit
		}
	}
	return nil
}

// checkErrorString checks if the error string is properly formatted.
// Returns the failure message and false if the string is not clean.
func checkErrorString(s string) (string, bool) {
	// Check for ending punctuation
	last, _ := utf8.DecodeLastRuneInString(s)
	if last == '.' || last == ':' || last == '!' || last == '\n' {
		return "error strings should not end with punctuation or a newline", false
	}

	first, firstN := utf8.DecodeRuneInString(s)
	if !unicode.IsUpper(first) {
		return "", true
	}

	// Allow proper nouns and acronyms (words with multiple uppercase letters or digits)
	for _, r := range s[firstN:] {
		if unicode.IsSpace(r) {
			break
		}
		if unicode.IsUpper(r) || unicode.IsDigit(r) {
			return "", true
		}
	}

	return "error strings should not be capitalized", false
}
