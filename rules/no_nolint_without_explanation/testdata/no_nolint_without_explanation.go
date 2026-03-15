package fixtures

// Invalid examples - nolint directives without explanations

//nolint:errcheck
// MATCH /nolint directive "//nolint:errcheck" should provide an explanation, e.g. //nolint:linter // reason/
func foo() {
	_ = doSomething()
}

//nolint:gosec
// MATCH /nolint directive "//nolint:gosec" should provide an explanation, e.g. //nolint:linter // reason/
func bar() {
	unsafeOperation()
}

//nolint:errcheck //
// MATCH /nolint directive "//nolint:errcheck //" should provide an explanation, e.g. //nolint:linter // reason/
func baz() {
	_ = riskyCall()
}

// Valid examples - nolint directives with explanations

//nolint:errcheck // error return is intentionally ignored here
func validFoo() {
	_ = doSomething()
}

//nolint:gosec // input is validated by the caller
func validBar() {
	unsafeOperation()
}

//nolint:errcheck,gosec // both checks are inapplicable in test code
func validBaz() {
	_ = riskyUncheckedCall()
}

// Bare //nolint without linters - handled by a different rule, not this one
//nolint
func ignoredBareNolint() {}

// Non-machine-readable format - handled by a different rule
// nolint:errcheck
func ignoredNonMachineReadable() {}

func doSomething() error        { return nil }
func unsafeOperation()          {}
func riskyCall() error          { return nil }
func riskyUncheckedCall() error { return nil }
