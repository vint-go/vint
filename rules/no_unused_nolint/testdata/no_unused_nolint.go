package fixtures

// Invalid examples - nolint directives that are unused

//nolint:errcheck
// MATCH /nolint directive "//nolint:errcheck" is unused/
func foo() int {
	return 42
}

//nolint:gosec
// MATCH /nolint directive "//nolint:gosec" is unused/
func bar() string {
	return "safe string"
}

//nolint
// MATCH /nolint directive "//nolint" is unused/
func baz() {
	println("hello")
}

//nolint:errcheck,gosec // both are suppressed but unused
// MATCH /nolint directive "//nolint:errcheck,gosec" is unused/
func qux() int {
	return 1
}

//nolint:errcheck // error is intentionally ignored
// MATCH /nolint directive "//nolint:errcheck" is unused/
func withExplanation() {
	println("still unused")
}

// Valid examples - no nolint directives present

// No nolint directive needed when there are no warnings
func validBar() int {
	return 42
}

// This is a normal comment mentioning nolinting in a sentence
func normalComment() {}

// nolintfoo is not a nolint directive
func notADirective() {}
