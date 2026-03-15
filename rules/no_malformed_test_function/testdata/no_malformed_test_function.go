package fixtures

import "testing"

// Invalid: test function has wrong parameter type.

func TestSomething(b *testing.B) { // MATCH /test function TestSomething has wrong signature: must be func TestSomething(t *testing.T)/
	_ = b
}

// Invalid: benchmark function has wrong parameter type.

func BenchmarkSomething(t *testing.T) { // MATCH /benchmark function BenchmarkSomething has wrong signature: must be func BenchmarkSomething(b *testing.B)/
	_ = t
}

// Invalid: fuzz function has wrong parameter type.

func FuzzSomething(t *testing.T) { // MATCH /fuzz function FuzzSomething has wrong signature: must be func FuzzSomething(f *testing.F)/
	_ = t
}

// Invalid: example function references non-existent identifier.

func ExampleNonExistentFunction() { // MATCH /example function ExampleNonExistentFunction references non-existent identifier NonExistentFunction/
	// Output: hello
}

// Valid: correct test function signature.

func TestGoodTest(t *testing.T) {
	_ = t
}

// Valid: correct benchmark function signature.

func BenchmarkGoodBench(b *testing.B) {
	_ = b
}

// Valid: correct fuzz function signature.

func FuzzGoodFuzz(f *testing.F) {
	_ = f
}

// Valid: example function for existing identifier.

func SomeFunction() {}

func ExampleSomeFunction() {
	// Output: hello
}

// Valid: package-level example (no identifier suffix).

func Example() {
	// Output: hello
}
