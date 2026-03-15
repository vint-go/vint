package fixtures

import "regexp"

// --- Invalid examples ---

func badDanglingAnchorStart() {
	// Anchor inside non-capturing group on only one alternative
	re := regexp.MustCompile(`(?:^foo|bar)`) // MATCH /dangling anchor ^ inside group: only some alternation branches are anchored in `^foo|bar`/
	_ = re
}

func badDanglingAnchorEnd() {
	re := regexp.MustCompile(`(?:foo|bar$)`) // MATCH /dangling anchor $ inside group: only some alternation branches are anchored in `foo|bar$`/
	_ = re
}

func badNestedQuantifiers() {
	re := regexp.MustCompile(`a**`) // MATCH /nested quantifier `*` applied to another quantifier/
	_ = re
}

func badNestedQuantifiersPlus() {
	re := regexp.MustCompile(`a*+`) // MATCH /nested quantifier `+` applied to another quantifier/
	_ = re
}

func badDuplicatedAlternation() {
	re := regexp.MustCompile(`a|b|a`) // MATCH /duplicated alternation `a`/
	_ = re
}

func badDuplicatedAlternationInGroup() {
	re := regexp.MustCompile(`(?:foo|bar|foo)`) // MATCH /duplicated alternation `foo`/
	_ = re
}

func badSuspiciousCharRange() {
	re := regexp.MustCompile(`[!-_]`) // MATCH /suspicious character range `!`-`_` in character class/
	_ = re
}

func badCharClassOverlap() {
	re := regexp.MustCompile(`[aba]`) // MATCH /duplicated character `a` in character class/
	_ = re
}

func badRedundantFlag() {
	re := regexp.MustCompile(`(?ii:foo)`) // MATCH /redundant flag `i` set multiple times/
	_ = re
}

func badConflictingFlag() {
	re := regexp.MustCompile(`(?i-i:foo)`) // MATCH /flag `i` is both set and cleared/
	_ = re
}

// --- Valid examples ---

func goodAnchorOutsideGroup() {
	re := regexp.MustCompile(`^(?:foo|bar)`)
	_ = re
}

func goodNoOverlappingChars() {
	re := regexp.MustCompile(`[ab]`)
	_ = re
}

func goodStandardRange() {
	re := regexp.MustCompile(`[a-z]`)
	_ = re
}

func goodDigitRange() {
	re := regexp.MustCompile(`[0-9]`)
	_ = re
}

func goodUppercaseRange() {
	re := regexp.MustCompile(`[A-Z]`)
	_ = re
}

func goodSimplePattern() {
	re := regexp.MustCompile(`^hello\s+world$`)
	_ = re
}

func goodAllAlternationsAnchored() {
	re := regexp.MustCompile(`(?:^foo|^bar)`)
	_ = re
}

func goodNonLiteralPattern() {
	pattern := `[abc]`
	re := regexp.MustCompile(pattern)
	_ = re
}

func goodCompile() {
	re, _ := regexp.Compile(`[a-z]+`)
	_ = re
}
