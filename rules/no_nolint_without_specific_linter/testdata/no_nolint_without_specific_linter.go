package fixtures

// Invalid examples - bare //nolint without specific linters

//nolint
// MATCH /nolint directive "//nolint" should mention specific linter(s), e.g. //nolint:errcheck/
func foo() {
	_ = doSomething()
}

//nolint // suppress all linters
// MATCH /nolint directive "//nolint // suppress all linters" should mention specific linter(s), e.g. //nolint:errcheck/
func bar() {
	unsafeOperation()
}

// Valid examples - nolint with specific linters

//nolint:errcheck
func validFoo() {
	_ = doSomething()
}

//nolint:gosec // this input is validated upstream
func validBar() {
	unsafeOperation()
}

//nolint:errcheck,gosec // both are acceptable here
func validBaz() {
	_ = riskyUncheckedCall()
}

func doSomething() error        { return nil }
func unsafeOperation()          {}
func riskyUncheckedCall() error { return nil }
