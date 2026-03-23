package no_unescaped_regexp_dot

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnescapedRegexpDotRule detects suspicious regexp patterns with unescaped dots
// before common domain extensions like com, org, net, info, etc.
type NoUnescapedRegexpDotRule struct{}

// Apply applies the rule to given file.
func (r *NoUnescapedRegexpDotRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnescapedRegexpDot{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnescapedRegexpDotRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintUnescapedRegexpDot{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnescapedRegexpDotRule) Name() string {
	return "noUnescapedRegexpDot"
}

// Group returns the rule group.
func (*NoUnescapedRegexpDotRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnescapedRegexpDotRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// commonDomainExtensions lists domain extensions that are commonly found
// after unescaped dots in regexp patterns.
var commonDomainExtensions = []string{
	"com", "org", "net", "info", "edu", "gov", "mil",
	"biz", "name", "mobi", "pro", "aero", "coop", "museum",
	"io", "co", "us", "uk", "ru", "de", "fr", "it", "es",
}

type lintUnescapedRegexpDot struct {
	onFailure func(lint.Failure)
}

func (w *lintUnescapedRegexpDot) Visit(node ast.Node) ast.Visitor {
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

	methodName := sel.Sel.Name
	if methodName != "Compile" && methodName != "MustCompile" &&
		methodName != "CompilePOSIX" && methodName != "MatchString" &&
		methodName != "Match" {
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

	if msg := checkUnescapedDomainDot(pattern); msg != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    msg,
		})
	}

	return w
}

// checkUnescapedDomainDot checks if a regexp pattern contains an unescaped dot
// followed by a common domain extension.
func checkUnescapedDomainDot(pattern string) string {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '\\' {
			// Skip escaped character
			i++
			continue
		}
		if pattern[i] == '[' {
			// Skip character class content
			i++
			if i < len(pattern) && pattern[i] == '^' {
				i++
			}
			if i < len(pattern) && pattern[i] == ']' {
				i++
			}
			for i < len(pattern) && pattern[i] != ']' {
				if pattern[i] == '\\' {
					i++
				}
				i++
			}
			continue
		}
		if pattern[i] == '.' {
			// Check if this dot is followed by a common domain extension
			rest := pattern[i+1:]
			for _, ext := range commonDomainExtensions {
				if strings.HasPrefix(rest, ext) {
					// Verify the extension is either at end or followed by non-alpha
					afterExt := rest[len(ext):]
					if afterExt == "" || !isAlpha(afterExt[0]) {
						return fmt.Sprintf("unescaped dot in regexp pattern before '%s', use '\\.' to match a literal dot", ext)
					}
				}
			}
		}
	}
	return ""
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
