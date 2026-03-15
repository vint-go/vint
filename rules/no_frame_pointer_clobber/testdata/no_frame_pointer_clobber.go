package fixtures

// Invalid: assembly function badExample clobbers BP without saving it first.
func badExample() int64 // MATCH /assembly function badExample clobbers frame pointer (BP) before saving it/

// Valid: assembly function goodExample saves and restores BP properly.
func goodExample() int64

// Valid: regular function with body, not an assembly stub.
func regularFunc(a, b int) int {
	return a + b
}
