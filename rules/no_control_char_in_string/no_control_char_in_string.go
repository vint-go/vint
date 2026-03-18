package no_control_char_in_string

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"unicode"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoControlCharInStringRule detects zero-width and control characters in string literals.
type NoControlCharInStringRule struct{}

// Apply applies the rule to given file.
func (r *NoControlCharInStringRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	ast.Inspect(file.AST, func(node ast.Node) bool {
		checkNode(node, file, onFailure)
		return true
	})

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoControlCharInStringRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	checkNode(node, file, onFailure)

	return failures
}

func checkNode(node ast.Node, file *lint.File, onFailure func(lint.Failure)) {
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	// Unquote the string to get the actual rune content
	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	for _, ch := range val {
		if isProblematicChar(ch) {
			onFailure(lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       lit,
				Failure:    fmt.Sprintf("string literal contains zero-width or control character U+%04X", ch),
			})
			return // report only the first problematic character per literal
		}
	}
}

// isProblematicChar returns true if the rune is a zero-width or control character
// that should not appear in string literals.
func isProblematicChar(r rune) bool {
	// Allow common whitespace control characters
	if r == '\n' || r == '\t' || r == '\r' {
		return false
	}

	// Check for Unicode control characters (Cc category)
	if unicode.Is(unicode.Cc, r) {
		return true
	}

	// Check for Unicode format characters (Cf category) with exceptions
	if unicode.Is(unicode.Cf, r) {
		// Allow specific Arabic characters
		if r >= 0x0600 && r <= 0x0605 {
			return false
		}
		if r == 0x0890 || r == 0x0891 {
			return false
		}
		if r == 0x08E2 {
			return false
		}

		// Allow variation selectors
		if unicode.Is(unicode.Variation_Selector, r) {
			return false
		}

		// Allow flag emoji tags (U+E0020-U+E007F)
		if r >= 0xE0020 && r <= 0xE007F {
			return false
		}

		return true
	}

	return false
}

// Name returns the rule name.
func (*NoControlCharInStringRule) Name() string {
	return "noControlCharInString"
}

// Group returns the rule group.
func (*NoControlCharInStringRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoControlCharInStringRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
