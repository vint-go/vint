package no_bad_regexp_pattern

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoBadRegexpPatternRule detects suspicious regexp patterns that may contain errors.
type NoBadRegexpPatternRule struct{}

// Apply applies the rule to given file.
func (r *NoBadRegexpPatternRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintBadRegexpPattern{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoBadRegexpPatternRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintBadRegexpPattern{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoBadRegexpPatternRule) Name() string {
	return "noBadRegexpPattern"
}

// Group returns the rule group.
func (*NoBadRegexpPatternRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoBadRegexpPatternRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBadRegexpPattern struct {
	onFailure func(lint.Failure)
}

func (w *lintBadRegexpPattern) Visit(node ast.Node) ast.Visitor {
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

	if sel.Sel.Name != "Compile" && sel.Sel.Name != "MustCompile" {
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

	issues := analyzeRegexpPattern(pattern)
	for _, issue := range issues {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    issue,
		})
	}

	return w
}

// analyzeRegexpPattern checks a regexp pattern for various issues.
func analyzeRegexpPattern(pattern string) []string {
	var issues []string

	if msgs := checkDanglingAnchors(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}
	if msgs := checkNestedQuantifiers(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}
	if msgs := checkDuplicatedAlternations(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}
	if msgs := checkSuspiciousCharRanges(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}
	if msgs := checkCharClassOverlaps(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}
	if msgs := checkRedundantFlags(pattern); len(msgs) > 0 {
		issues = append(issues, msgs...)
	}

	return issues
}

// checkDanglingAnchors detects ^ or $ positioned improperly in alternation groups.
// For example: (?:^foo|bar) has ^ only on the first alternative.
func checkDanglingAnchors(pattern string) []string {
	var issues []string

	// Look for groups containing alternations where anchors appear only on some branches.
	// We search for patterns like (?:^...|...) or (?:...|...$)
	groups := extractGroups(pattern)
	for _, g := range groups {
		if !strings.Contains(g.content, "|") {
			continue
		}
		alts := splitAlternations(g.content)
		if len(alts) < 2 {
			continue
		}

		hasAnchorStart := make([]bool, len(alts))
		hasAnchorEnd := make([]bool, len(alts))
		anyStart := false
		anyEnd := false
		allStart := true
		allEnd := true

		for i, alt := range alts {
			if strings.HasPrefix(alt, "^") {
				hasAnchorStart[i] = true
				anyStart = true
			} else {
				allStart = false
			}
			if strings.HasSuffix(alt, "$") && !strings.HasSuffix(alt, "\\$") {
				hasAnchorEnd[i] = true
				anyEnd = true
			} else {
				allEnd = false
			}
		}

		if anyStart && !allStart {
			issues = append(issues, fmt.Sprintf("dangling anchor ^ inside group: only some alternation branches are anchored in `%s`", g.content))
		}
		if anyEnd && !allEnd {
			issues = append(issues, fmt.Sprintf("dangling anchor $ inside group: only some alternation branches are anchored in `%s`", g.content))
		}
	}

	return issues
}

// group represents a group found in a regexp pattern.
type group struct {
	content string // content inside the group (without parens/prefix)
}

// extractGroups extracts groups from a regexp pattern.
func extractGroups(pattern string) []group {
	var groups []group
	i := 0
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, size2 := utf8.DecodeRuneInString(pattern[i:])
				i += size2
			}
			continue
		}
		if r == '(' {
			// Find the matching close paren
			start := i + size
			// Skip group prefix like ?: ?= etc.
			if start < len(pattern) && pattern[start] == '?' {
				start++
				for start < len(pattern) {
					c, sz := utf8.DecodeRuneInString(pattern[start:])
					if c == ':' || c == '=' || c == '!' || c == '<' || c == '>' {
						start += sz
						break
					}
					if c == ')' {
						break
					}
					// Flag characters
					if c == 'i' || c == 'm' || c == 's' || c == 'U' || c == '-' {
						start += sz
						continue
					}
					break
				}
			}
			depth := 1
			end := i + size
			for end < len(pattern) && depth > 0 {
				c, sz := utf8.DecodeRuneInString(pattern[end:])
				if c == '\\' {
					end += sz
					if end < len(pattern) {
						_, sz2 := utf8.DecodeRuneInString(pattern[end:])
						end += sz2
					}
					continue
				}
				if c == '(' {
					depth++
				} else if c == ')' {
					depth--
				}
				if depth > 0 {
					end += sz
				} else {
					break
				}
			}
			if depth == 0 && start <= end {
				content := pattern[start:end]
				groups = append(groups, group{content: content})
			}
			i = end + 1
			continue
		}
		i += size
	}
	return groups
}

// splitAlternations splits a regex string by top-level | (not inside nested groups or char classes).
func splitAlternations(s string) []string {
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
				_, size2 := utf8.DecodeRuneInString(s[i:])
				i += size2
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

// checkNestedQuantifiers detects repeated greedy quantifiers like a** or a*+.
func checkNestedQuantifiers(pattern string) []string {
	var issues []string
	quantifiers := map[rune]bool{'*': true, '+': true, '?': true}

	i := 0
	prevWasQuantifier := false
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, size2 := utf8.DecodeRuneInString(pattern[i:])
				i += size2
			}
			prevWasQuantifier = false
			continue
		}
		if r == '[' {
			// skip char class
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
				if c == ']' {
					i += sz
					break
				}
				i += sz
			}
			prevWasQuantifier = false
			continue
		}
		if quantifiers[r] {
			if prevWasQuantifier {
				issues = append(issues, fmt.Sprintf("nested quantifier `%c` applied to another quantifier", r))
			}
			prevWasQuantifier = true
		} else if r == '{' {
			// {n,m} style quantifier
			j := i + size
			for j < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[j:])
				if c == '}' {
					j += sz
					break
				}
				j += sz
			}
			if prevWasQuantifier {
				issues = append(issues, "nested quantifier `{` applied to another quantifier")
			}
			prevWasQuantifier = true
			i = j
			continue
		} else {
			prevWasQuantifier = false
		}
		i += size
	}
	return issues
}

// checkDuplicatedAlternations finds identical options in alternation groups.
func checkDuplicatedAlternations(pattern string) []string {
	var issues []string

	// Check the top-level pattern
	if strings.Contains(pattern, "|") {
		alts := splitAlternations(pattern)
		if dupes := findDuplicates(alts); len(dupes) > 0 {
			for _, d := range dupes {
				issues = append(issues, fmt.Sprintf("duplicated alternation `%s`", d))
			}
		}
	}

	// Check inside groups
	groups := extractGroups(pattern)
	for _, g := range groups {
		if !strings.Contains(g.content, "|") {
			continue
		}
		alts := splitAlternations(g.content)
		if dupes := findDuplicates(alts); len(dupes) > 0 {
			for _, d := range dupes {
				issues = append(issues, fmt.Sprintf("duplicated alternation `%s`", d))
			}
		}
	}

	return issues
}

func findDuplicates(items []string) []string {
	seen := map[string]int{}
	var dupes []string
	for _, item := range items {
		seen[item]++
		if seen[item] == 2 {
			dupes = append(dupes, item)
		}
	}
	return dupes
}

// checkSuspiciousCharRanges flags non-standard ranges in character classes.
func checkSuspiciousCharRanges(pattern string) []string {
	var issues []string

	i := 0
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, size2 := utf8.DecodeRuneInString(pattern[i:])
				i += size2
			}
			continue
		}
		if r == '[' {
			i += size
			// Check if negated
			if i < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[i:])
				if c == '^' {
					i += sz
				}
			}
			// Parse the character class content
			classContent := parseCharClassContent(pattern, &i)
			if msgs := analyzeCharClassRanges(classContent); len(msgs) > 0 {
				issues = append(issues, msgs...)
			}
			continue
		}
		i += size
	}

	return issues
}

// charClassElement represents an element in a character class.
type charClassElement struct {
	isRange bool
	lo      rune
	hi      rune
	single  rune
}

// parseCharClassContent parses the elements of a character class and advances i past the closing ].
func parseCharClassContent(pattern string, i *int) []charClassElement {
	var elements []charClassElement
	// Allow ] as first char in class
	first := true

	for *i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[*i:])

		if r == ']' && !first {
			*i += size
			return elements
		}
		first = false

		if r == '\\' {
			*i += size
			if *i < len(pattern) {
				escaped, sz := utf8.DecodeRuneInString(pattern[*i:])
				*i += sz
				// Check for range
				if *i < len(pattern) && pattern[*i] == '-' && *i+1 < len(pattern) && pattern[*i+1] != ']' {
					*i++ // skip -
					hi, hiSz := readCharClassChar(pattern, i)
					elements = append(elements, charClassElement{isRange: true, lo: escaped, hi: hi})
					_ = hiSz
				} else {
					elements = append(elements, charClassElement{single: escaped})
				}
			}
			continue
		}

		*i += size
		// Check for range: c-d
		if *i < len(pattern) && pattern[*i] == '-' && *i+1 < len(pattern) && pattern[*i+1] != ']' {
			*i++ // skip -
			hi, _ := readCharClassChar(pattern, i)
			elements = append(elements, charClassElement{isRange: true, lo: r, hi: hi})
		} else {
			elements = append(elements, charClassElement{single: r})
		}
	}
	return elements
}

func readCharClassChar(pattern string, i *int) (rune, int) {
	if *i >= len(pattern) {
		return 0, 0
	}
	r, size := utf8.DecodeRuneInString(pattern[*i:])
	if r == '\\' {
		*i += size
		if *i < len(pattern) {
			r2, sz2 := utf8.DecodeRuneInString(pattern[*i:])
			*i += sz2
			return r2, sz2
		}
		return r, size
	}
	*i += size
	return r, size
}

// analyzeCharClassRanges checks for suspicious ranges in character class elements.
func analyzeCharClassRanges(elements []charClassElement) []string {
	var issues []string

	for _, elem := range elements {
		if !elem.isRange {
			continue
		}
		if !isStandardRange(elem.lo, elem.hi) {
			issues = append(issues, fmt.Sprintf("suspicious character range `%c`-`%c` in character class", elem.lo, elem.hi))
		}
	}

	return issues
}

// isStandardRange checks if a character range is a standard one (letter-letter or digit-digit).
func isStandardRange(lo, hi rune) bool {
	if lo > hi {
		return false // invalid range
	}
	// digit-digit
	if lo >= '0' && lo <= '9' && hi >= '0' && hi <= '9' {
		return true
	}
	// lowercase letter to lowercase letter
	if unicode.IsLower(lo) && unicode.IsLower(hi) {
		return true
	}
	// uppercase letter to uppercase letter
	if unicode.IsUpper(lo) && unicode.IsUpper(hi) {
		return true
	}
	return false
}

// checkCharClassOverlaps detects duplicate or intersecting character definitions within classes.
func checkCharClassOverlaps(pattern string) []string {
	var issues []string

	i := 0
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, size2 := utf8.DecodeRuneInString(pattern[i:])
				i += size2
			}
			continue
		}
		if r == '[' {
			i += size
			// Check if negated
			if i < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[i:])
				if c == '^' {
					i += sz
				}
			}

			classContent := parseCharClassContent(pattern, &i)
			if msgs := findCharClassDuplicates(classContent); len(msgs) > 0 {
				issues = append(issues, msgs...)
			}
			continue
		}
		i += size
	}

	return issues
}

// findCharClassDuplicates detects duplicate single characters in a character class.
func findCharClassDuplicates(elements []charClassElement) []string {
	var issues []string
	seen := map[rune]int{}

	for _, elem := range elements {
		if !elem.isRange {
			seen[elem.single]++
			if seen[elem.single] == 2 {
				issues = append(issues, fmt.Sprintf("duplicated character `%c` in character class", elem.single))
			}
		}
	}

	return issues
}

// checkRedundantFlags reports setting flags twice or clearing unset flags.
func checkRedundantFlags(pattern string) []string {
	var issues []string

	i := 0
	for i < len(pattern) {
		r, size := utf8.DecodeRuneInString(pattern[i:])
		if r == '\\' {
			i += size
			if i < len(pattern) {
				_, size2 := utf8.DecodeRuneInString(pattern[i:])
				i += size2
			}
			continue
		}
		if r == '(' && i+size < len(pattern) && pattern[i+size] == '?' {
			// Parse flag group (?flags) or (?flags:...)
			flagStart := i + size + 1 // skip (?
			flagEnd := flagStart
			clearing := false
			setFlags := map[rune]int{}
			clearFlags := map[rune]int{}

			for flagEnd < len(pattern) {
				c, sz := utf8.DecodeRuneInString(pattern[flagEnd:])
				if c == ':' || c == ')' {
					break
				}
				if c == '-' {
					clearing = true
					flagEnd += sz
					continue
				}
				if isRegexpFlag(c) {
					if clearing {
						clearFlags[c]++
						if clearFlags[c] == 2 {
							issues = append(issues, fmt.Sprintf("redundant flag `%c` cleared multiple times", c))
						}
					} else {
						setFlags[c]++
						if setFlags[c] == 2 {
							issues = append(issues, fmt.Sprintf("redundant flag `%c` set multiple times", c))
						}
					}
				}
				flagEnd += sz
			}

			// Check for flags that are both set and cleared
			for flag := range clearFlags {
				if setFlags[flag] > 0 {
					issues = append(issues, fmt.Sprintf("flag `%c` is both set and cleared", flag))
				}
			}

			i = flagEnd
			continue
		}
		i += size
	}

	return issues
}

func isRegexpFlag(r rune) bool {
	return r == 'i' || r == 'm' || r == 's' || r == 'U'
}
