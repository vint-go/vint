package fixtures

// Invalid: unexported function never called or referenced from any entry point.

func unusedHelper() int { // MATCH /func unusedHelper is unused/
	return 42
}

func anotherUnused() string { // MATCH /func anotherUnused is unused/
	return "hello"
}

// Invalid: mutually recursive but never called from an entry point.

func mutualA() int { // MATCH /func mutualA is unused/
	return mutualB()
}

func mutualB() int { // MATCH /func mutualB is unused/
	return mutualA()
}

// Valid: exported functions are considered used by the package.

func PublicAPI() string {
	return internalHelper()
}

func internalHelper() string {
	return "data"
}

// Valid: function is called by an exported function (transitive reachability).

func ExportedCaller() string {
	return callee()
}

func callee() string {
	return "result"
}

// Valid: function used as a value, reachable from exported function.

func ExportedUseFunc() int {
	return apply(myFunc)
}

func apply(f func() int) int {
	return f()
}

func myFunc() int {
	return 1
}
