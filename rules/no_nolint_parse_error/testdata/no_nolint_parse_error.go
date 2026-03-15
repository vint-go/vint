package fixtures

// Invalid examples - should trigger failures

//nolint:errcheck gosec
// MATCH /malformed nolint directive "//nolint:errcheck gosec": expected format is //nolint[:<linter1>,<linter2>] [// explanation]/
func foo() {}

//nolint: errcheck, ,gosec
// MATCH /malformed nolint directive "//nolint: errcheck, ,gosec": expected format is //nolint[:<linter1>,<linter2>] [// explanation]/
func bar() {}

//nolint:errcheck:gosec
// MATCH /malformed nolint directive "//nolint:errcheck:gosec": expected format is //nolint[:<linter1>,<linter2>] [// explanation]/
func baz() {}

// Valid examples - should NOT trigger failures

//nolint:errcheck
func validFoo() {}

//nolint:errcheck,gosec
func validBar() {}

//nolint:errcheck,gosec // both are safe in this context
func validBaz() {}

//nolint
func validBare() {}

//nolint // some reason
func validBareWithReason() {}

// This is a normal comment
func normalComment() {}
