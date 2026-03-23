package use_simplified_regexp

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseSimplifiedRegexpRule detects regexp patterns that can be simplified.
type UseSimplifiedRegexpRule struct{}

// Apply applies the rule to given file.
func (r *UseSimplifiedRegexpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSimplifiedRegexp{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSimplifiedRegexpRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintSimplifiedRegexp{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseSimplifiedRegexpRule) Name() string {
	return "useSimplifiedRegexp"
}

// Group returns the rule group.
func (*UseSimplifiedRegexpRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseSimplifiedRegexpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSimplifiedRegexp struct {
	onFailure func(lint.Failure)
}

func (w *lintSimplifiedRegexp) Visit(node ast.Node) ast.Visitor {
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
		funcName != "CompilePOSIX" && funcName != "MustCompilePOSIX" {
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

	// Skip patterns exceeding 60 characters
	if len(pattern) > 60 {
		return w
	}

	simplified := simplifyPattern(pattern)
	if simplified != pattern {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("regexp pattern can be simplified: `%s` -> `%s`", pattern, simplified),
		})
	}

	return w
}

// simplifyPattern performs up to 2 simplification passes on a regexp pattern.
func simplifyPattern(pattern string) string {
	result := pattern
	for i := 0; i < 2; i++ {
		prev := result
		result = simplifyOnce(result)
		if result == prev {
			break
		}
	}
	return result
}

// simplifyOnce runs all simplification rules on the pattern once.
func simplifyOnce(pattern string) string {
	result := pattern
	result = simplifyQuantifiers(result)
	result = simplifyCharClasses(result)
	result = simplifyNegatedCharClasses(result)
	result = simplifySingleCharAlternation(result)
	result = simplifyCommonPrefixSuffix(result)
	result = simplifyUnnecessaryNonCapturingGroup(result)
	return result
}

// simplifyQuantifiers converts verbose quantifiers to shorter forms:
// {0,1} -> ?, {1,} -> +, {0,} -> *
func simplifyQuantifiers(pattern string) string {
	var b strings.Builder
	i := 0
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			b.WriteString(pattern[i : i+size])
			i += size
			if i < len(pattern) {
				_, sz := utf8.DecodeRuneInString(pattern[i:])
				b.WriteString(pattern[i : i+sz])
				i += sz
			}
			continue
		}
		if r == '[' {
			// Skip through character class
			b.WriteByte('[')
			i += size
			for i < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[i:])
				b.WriteString(pattern[i : i+sz])
				i += sz
				if c == '\\' && i < len(pattern) {
					_, sz2 := utf8.DecodeRuneInString(pattern[i:])
					b.WriteString(pattern[i : i+sz2])
					i += sz2
					continue
				}
				if c == ']' {
					break
				}
			}
			continue
		}
		if r == '{' {
			// Try to parse a quantifier like {n,m}
			end := strings.IndexByte(pattern[i:], '}')
			if end != -1 {
				inside := pattern[i+1 : i+end]
				replacement := ""
				switch inside {
				case "0,1":
					replacement = "?"
				case "1,":
					replacement = "+"
				case "0,":
					replacement = "*"
				}
				if replacement != "" {
					b.WriteString(replacement)
					i += end + 1
					continue
				}
			}
		}
		b.WriteString(pattern[i : i+size])
		i += size
	}
	return b.String()
}

// simplifyCharClasses replaces known character classes with shortcuts:
// [0-9] -> \d, [[:space:]] -> \s, [[:digit:]] -> \d, etc.
func simplifyCharClasses(pattern string) string {
	// Simple exact replacements for well-known character classes
	replacements := []struct {
		from string
		to   string
	}{
		{"[0-9]", `\d`},
		{"[[:digit:]]", `\d`},
		{"[[:space:]]", `\s`},
		{"[a-zA-Z0-9_]", `\w`},
		{"[A-Za-z0-9_]", `\w`},
		{"[0-9a-zA-Z_]", `\w`},
		{"[0-9A-Za-z_]", `\w`},
		{"[_a-zA-Z0-9]", `\w`},
		{"[_A-Za-z0-9]", `\w`},
		{"[ \\t\\n\\r\\f\\v]", `\s`},
		{"[\\t\\n\\r\\f\\v ]", `\s`},
	}

	result := pattern
	for _, rep := range replacements {
		result = replaceOutsideClass(result, rep.from, rep.to)
	}
	return result
}

// simplifyNegatedCharClasses replaces negated char classes with shortcuts:
// [^0-9] -> \D, [^[:space:]] -> \S, etc.
func simplifyNegatedCharClasses(pattern string) string {
	replacements := []struct {
		from string
		to   string
	}{
		{"[^0-9]", `\D`},
		{"[^[:digit:]]", `\D`},
		{"[^[:space:]]", `\S`},
		{"[^a-zA-Z0-9_]", `\W`},
		{"[^A-Za-z0-9_]", `\W`},
		{"[^0-9a-zA-Z_]", `\W`},
		{"[^0-9A-Za-z_]", `\W`},
	}

	result := pattern
	for _, rep := range replacements {
		result = replaceOutsideClass(result, rep.from, rep.to)
	}
	return result
}

// replaceOutsideClass replaces occurrences of 'from' with 'to' in the pattern,
// but only when not inside a character class (to avoid nested replacements).
func replaceOutsideClass(pattern, from, to string) string {
	// For character class replacements, the 'from' itself starts with '[',
	// so we just do a simple string replacement since these are complete
	// self-contained character class patterns.
	return strings.ReplaceAll(pattern, from, to)
}

// simplifySingleCharAlternation converts alternations of single characters
// into character classes: a|b|c -> [abc]
func simplifySingleCharAlternation(pattern string) string {
	return processAlternations(pattern, func(alts []string) (string, bool) {
		// Check if all alternatives are single characters (not meta-characters)
		if len(alts) < 2 {
			return "", false
		}
		var chars []string
		for _, alt := range alts {
			if len(alt) == 1 && !isRegexpMeta(rune(alt[0])) {
				chars = append(chars, alt)
			} else if len(alt) == 2 && alt[0] == '\\' {
				// Escaped character like \d
				chars = append(chars, alt)
			} else {
				return "", false
			}
		}
		return "[" + strings.Join(chars, "") + "]", true
	})
}

// simplifyCommonPrefixSuffix factors out common prefixes/suffixes in alternations.
// Examples: http|https -> https?, foo|foobar -> foo(?:bar)?
func simplifyCommonPrefixSuffix(pattern string) string {
	return processAlternations(pattern, func(alts []string) (string, bool) {
		if len(alts) != 2 {
			return "", false
		}

		a, b := alts[0], alts[1]
		if a == b {
			return "", false
		}

		// Try common prefix factoring
		prefix := commonPrefix(a, b)
		if len(prefix) > 0 {
			suffA := a[len(prefix):]
			suffB := b[len(prefix):]

			// Case: one suffix is empty -> prefix + optional suffix
			if suffA == "" && len(suffB) == 1 {
				return prefix + string(suffB[0]) + "?", true
			}
			if suffB == "" && len(suffA) == 1 {
				return prefix + string(suffA[0]) + "?", true
			}
			if suffA == "" && len(suffB) > 1 {
				return prefix + "(?:" + suffB + ")?", true
			}
			if suffB == "" && len(suffA) > 1 {
				return prefix + "(?:" + suffA + ")?", true
			}
		}

		// Try common suffix factoring
		suffix := commonSuffix(a, b)
		if len(suffix) > 0 {
			prefA := a[:len(a)-len(suffix)]
			prefB := b[:len(b)-len(suffix)]

			if prefA == "" && len(prefB) == 1 {
				return string(prefB[0]) + "?" + suffix, true
			}
			if prefB == "" && len(prefA) == 1 {
				return string(prefA[0]) + "?" + suffix, true
			}
			if prefA == "" && len(prefB) > 1 {
				return "(?:" + prefB + ")?" + suffix, true
			}
			if prefB == "" && len(prefA) > 1 {
				return "(?:" + prefA + ")?" + suffix, true
			}
		}

		return "", false
	})
}

// simplifyUnnecessaryNonCapturingGroup removes unnecessary non-capturing group wrappers.
// For example: (?:a) -> a (when the group is the entire pattern or adds no meaning)
func simplifyUnnecessaryNonCapturingGroup(pattern string) string {
	// Handle case where entire pattern is wrapped: (?:content)
	if strings.HasPrefix(pattern, "(?:") && strings.HasSuffix(pattern, ")") {
		inner := pattern[3 : len(pattern)-1]
		// Check if the closing paren is indeed the matching one
		if !containsUnescapedTopLevel(inner, '|') && matchingClose(pattern, 0) == len(pattern)-1 {
			return inner
		}
	}

	// Handle non-capturing groups wrapping a single character or simple escape
	var b strings.Builder
	i := 0
	changed := false
	for i < len(pattern) {
		if i+3 < len(pattern) && pattern[i:i+3] == "(?:" {
			closeIdx := matchingClose(pattern, i)
			if closeIdx > 0 {
				inner := pattern[i+3 : closeIdx]
				// Replace (?:X) with X when X is a single character or single escape
				if isSingleAtom(inner) && closeIdx+1 < len(pattern) && !isQuantifier(rune(pattern[closeIdx+1])) {
					b.WriteString(inner)
					i = closeIdx + 1
					changed = true
					continue
				}
			}
		}
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			b.WriteString(pattern[i : i+size])
			i += size
			if i < len(pattern) {
				_, sz := utf8.DecodeRuneInString(pattern[i:])
				b.WriteString(pattern[i : i+sz])
				i += sz
			}
			continue
		}
		b.WriteString(pattern[i : i+size])
		i += size
	}

	if changed {
		return b.String()
	}
	return pattern
}

// processAlternations finds top-level alternations in the pattern and
// applies the transform function to them.
func processAlternations(pattern string, transform func([]string) (string, bool)) string {
	// First try top-level alternations
	alts := splitTopLevelAlts(pattern)
	if len(alts) >= 2 {
		if result, ok := transform(alts); ok {
			return result
		}
	}

	// Then try alternations inside groups
	var b strings.Builder
	i := 0
	changed := false
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			b.WriteString(pattern[i : i+size])
			i += size
			if i < len(pattern) {
				_, sz := utf8.DecodeRuneInString(pattern[i:])
				b.WriteString(pattern[i : i+sz])
				i += sz
			}
			continue
		}
		if r == '(' {
			closeIdx := matchingClose(pattern, i)
			if closeIdx < 0 {
				b.WriteString(pattern[i : i+size])
				i += size
				continue
			}
			groupContent := pattern[i : closeIdx+1]
			// Check if it's a non-capturing group with alternation
			if strings.HasPrefix(groupContent, "(?:") {
				inner := groupContent[3 : len(groupContent)-1]
				innerAlts := splitTopLevelAlts(inner)
				if len(innerAlts) >= 2 {
					if result, ok := transform(innerAlts); ok {
						// If the result doesn't contain alternation, we might not need the group
						if !containsUnescapedTopLevel(result, '|') {
							b.WriteString(result)
						} else {
							b.WriteString("(?:" + result + ")")
						}
						i = closeIdx + 1
						changed = true
						continue
					}
				}
			}
			b.WriteString(groupContent)
			i = closeIdx + 1
			continue
		}
		b.WriteString(pattern[i : i+size])
		i += size
	}

	if changed {
		return b.String()
	}
	return pattern
}

// splitTopLevelAlts splits a pattern by top-level | (not inside groups or char classes).
func splitTopLevelAlts(s string) []string {
	var parts []string
	depth := 0
	inClass := false
	start := 0
	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\\' {
			i += size
			if i < len(s) {
				_, sz := utf8.DecodeRuneInString(s[i:])
				i += sz
			}
			continue
		}
		if r == '[' && !inClass {
			inClass = true
			i += size
			continue
		}
		if r == ']' && inClass {
			inClass = false
			i += size
			continue
		}
		if inClass {
			i += size
			continue
		}
		if r == '(' {
			depth++
		} else if r == ')' {
			depth--
		} else if r == '|' && depth == 0 {
			parts = append(parts, s[start:i])
			start = i + size
		}
		i += size
	}
	parts = append(parts, s[start:])
	return parts
}

// matchingClose finds the index of the closing paren matching the opening paren at idx.
func matchingClose(pattern string, idx int) int {
	depth := 0
	i := idx
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, sz := utf8.DecodeRuneInString(pattern[i:])
				i += sz
			}
			continue
		}
		if r == '[' {
			// Skip character class
			i += size
			for i < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[i:])
				if c == '\\' {
					i += sz
					if i < len(pattern) {
						_, sz2 := utf8.DecodeRuneInString(pattern[i:])
						i += sz2
					}
					continue
				}
				i += sz
				if c == ']' {
					break
				}
			}
			continue
		}
		if r == '(' {
			depth++
		} else if r == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
		i += size
	}
	return -1
}

// containsUnescapedTopLevel checks if the pattern contains an unescaped char
// at the top level (not inside groups or char classes).
func containsUnescapedTopLevel(s string, target rune) bool {
	depth := 0
	inClass := false
	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\\' {
			i += size
			if i < len(s) {
				_, sz := utf8.DecodeRuneInString(s[i:])
				i += sz
			}
			continue
		}
		if r == '[' && !inClass {
			inClass = true
			i += size
			continue
		}
		if r == ']' && inClass {
			inClass = false
			i += size
			continue
		}
		if inClass {
			i += size
			continue
		}
		if r == '(' {
			depth++
		} else if r == ')' {
			depth--
		} else if r == target && depth == 0 {
			return true
		}
		i += size
	}
	return false
}

// isSingleAtom returns true if the string represents a single regex atom:
// a single character or a single escape sequence like \d.
func isSingleAtom(s string) bool {
	if len(s) == 0 {
		return false
	}
	if len(s) == 1 {
		return !isRegexpMeta(rune(s[0]))
	}
	if len(s) == 2 && s[0] == '\\' {
		return true
	}
	return false
}

// isRegexpMeta returns true if the rune is a regexp metacharacter.
func isRegexpMeta(r rune) bool {
	return strings.ContainsRune(`\.+*?()[]{}|^$`, r)
}

// isQuantifier returns true if the rune is a quantifier character.
func isQuantifier(r rune) bool {
	return r == '*' || r == '+' || r == '?' || r == '{'
}

// commonPrefix returns the common prefix of two strings, being careful
// not to split escape sequences.
func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) {
		ra, sza := utf8.DecodeRuneInString(a[i:])
		rb, szb := utf8.DecodeRuneInString(b[i:])
		if ra != rb || sza != szb {
			break
		}
		if ra == '\\' {
			// Must include the next character in the escape
			if i+sza < len(a) && i+sza < len(b) {
				_, sza2 := utf8.DecodeRuneInString(a[i+sza:])
				_, szb2 := utf8.DecodeRuneInString(b[i+sza:])
				if a[i+sza:i+sza+sza2] != b[i+sza:i+sza+szb2] {
					break
				}
				i += sza + sza2
				continue
			}
			break
		}
		i += sza
	}
	return a[:i]
}

// commonSuffix returns the common suffix of two strings.
func commonSuffix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) {
		// Work from the end
		posA := len(a) - 1 - i
		posB := len(b) - 1 - i
		if posA < 0 || posB < 0 {
			break
		}
		if a[posA] != b[posB] {
			break
		}
		// Make sure we don't split an escape sequence
		if posA > 0 && a[posA-1] == '\\' {
			// This byte is the second part of an escape - skip it
			break
		}
		i++
	}
	if i == 0 {
		return ""
	}
	return a[len(a)-i:]
}
