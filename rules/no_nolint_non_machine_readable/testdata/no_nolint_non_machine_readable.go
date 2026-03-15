package fixtures

// Invalid examples - should trigger failures

// nolint:errcheck
// MATCH /nolint directive "// nolint:errcheck" is not machine-readable, use "//nolint:errcheck"/
func foo() {}

// nolint:gosec // reason for suppression
// MATCH /nolint directive "// nolint:gosec // reason for suppression" is not machine-readable, use "//nolint:gosec // reason for suppression"/
func bar() {}

// nolint
// MATCH /nolint directive "// nolint" is not machine-readable, use "//nolint"/
func baz() {}

// Valid examples - should NOT trigger failures

//nolint:errcheck
func validFoo() {}

//nolint:gosec // this is safe in this context
func validBar() {}

//nolint:errcheck,gosec // multiple linters suppressed for valid reason
func validBaz() {}

// This is a normal comment, not a nolint directive
func normalComment() {}
