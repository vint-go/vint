package use_modern_octal_literal

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseModernOctalLiteralRule detects old-style octal literals (e.g. 0755) and suggests
// using the modern 0o prefix style introduced in Go 1.13.
type UseModernOctalLiteralRule struct{}

// Apply applies the rule to given file.
func (r *UseModernOctalLiteralRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintModernOctal{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseModernOctalLiteralRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintModernOctal{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseModernOctalLiteralRule) Name() string {
	return "useModernOctalLiteral"
}

// Group returns the rule group.
func (*UseModernOctalLiteralRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseModernOctalLiteralRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintModernOctal struct {
	onFailure func(lint.Failure)
}

func (w *lintModernOctal) Visit(node ast.Node) ast.Visitor {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind != token.INT {
		return w
	}

	value := lit.Value

	// Must start with '0' and have at least 2 characters
	if len(value) < 2 {
		return w
	}

	// Must start with '0' but NOT be a hex (0x/0X), binary (0b/0B), or modern octal (0o/0O) literal
	if value[0] != '0' {
		return w
	}

	// Skip if it's just "0"
	if len(value) == 1 {
		return w
	}

	second := value[1]

	// Skip hex, binary, modern octal prefixes
	if second == 'x' || second == 'X' || second == 'b' || second == 'B' || second == 'o' || second == 'O' {
		return w
	}

	// At this point, it starts with '0' followed by something else.
	// Check that the remaining characters are valid octal digits (0-7) or underscores.
	// If any digit is 8 or 9, it's not a valid old-style octal literal (it's a decimal with leading zero
	// which would be a compile error in Go, but we skip it).
	isOldOctal := false
	for i := 1; i < len(value); i++ {
		ch := value[i]
		if ch == '_' {
			continue
		}
		if ch >= '0' && ch <= '7' {
			isOldOctal = true
			continue
		}
		// If we see 8 or 9, or any non-digit, it's not an old-style octal
		return w
	}

	if !isOldOctal {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       lit,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("old-style octal literal %s, use %s instead", value, modernOctal(value)),
	})

	return w
}

// modernOctal converts an old-style octal literal to the modern 0o prefix style.
func modernOctal(value string) string {
	// Remove leading '0' and add '0o' prefix
	return "0o" + value[1:]
}
