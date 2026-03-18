package fixtures

import "regexp"

func invalidRegexps() {
	// Invalid regex - unclosed group
	re := regexp.MustCompile("foo(bar") // MATCH /invalid regular expression: error parsing regexp: missing closing ): `foo(bar`/
	_ = re

	// Invalid regex - unclosed character class
	re = regexp.MustCompile("[abc") // MATCH /invalid regular expression: error parsing regexp: missing closing ]: `[abc`/

	// Invalid regex - unclosed group in Compile
	_, _ = regexp.Compile("(abc") // MATCH /invalid regular expression: error parsing regexp: missing closing ): `(abc`/

	// Invalid regex - bad escape
	re = regexp.MustCompile(`\`) // MATCH /invalid regular expression: error parsing regexp: trailing backslash at end of expression: ``/

	// Invalid regex in Match
	_, _ = regexp.Match("(", nil) // MATCH /invalid regular expression: error parsing regexp: missing closing ): `(`/

	// Invalid regex in MatchString
	_, _ = regexp.MatchString("[", "") // MATCH /invalid regular expression: error parsing regexp: missing closing ]: `[`/

	// Invalid regex in CompilePOSIX
	_, _ = regexp.CompilePOSIX("(abc") // MATCH /invalid regular expression: error parsing regexp: missing closing ): `(abc`/

	// Invalid regex - bad repetition with nothing to repeat
	re = regexp.MustCompile("*") // MATCH /invalid regular expression: error parsing regexp: missing argument to repetition operator: `*`/
}

func validRegexps() {
	// Valid regex
	re := regexp.MustCompile("foo(bar)")
	_ = re

	// Valid regex - character class
	re = regexp.MustCompile("[abc]")

	// Valid regex - repetition
	re = regexp.MustCompile("foo{3}")

	// Valid regex - complex pattern
	re = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	// Valid regex - empty pattern
	re = regexp.MustCompile("")

	// Valid regex - alternation
	re = regexp.MustCompile("foo|bar")

	// Valid regex - special characters escaped
	re = regexp.MustCompile(`foo\.bar`)

	// Non-literal pattern (should be skipped)
	pattern := "foo(bar"
	re = regexp.MustCompile(pattern)

	// Valid regex via Compile
	_, _ = regexp.Compile("^[a-z]+$")

	// Valid regex via Match
	_, _ = regexp.Match("^[a-z]+$", nil)

	// Valid regex via MatchString
	_, _ = regexp.MatchString("^[a-z]+$", "")

	// Valid regex via CompilePOSIX
	_, _ = regexp.CompilePOSIX("^[a-z]+$")

	// foo{ is valid in Go regexp (literal brace)
	re = regexp.MustCompile("foo{")
}
