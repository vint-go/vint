package fixtures

import "testing"

// Invalid: assigning to b.N in a benchmark function.
func BenchmarkBad(b *testing.B) {
	b.N = 1000 // MATCH /assignment to b.N in benchmark distorts the results/
	for i := 0; i < b.N; i++ {
		// do work
	}
}

// Invalid: assigning to b.N with += operator.
func BenchmarkBadCompound(b *testing.B) {
	b.N += 100 // MATCH /assignment to b.N in benchmark distorts the results/
	for i := 0; i < b.N; i++ {
		// do work
	}
}

// Valid: reading b.N is fine.
func BenchmarkGood(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// do work
	}
}

// Valid: not a benchmark function.
func TestNotBenchmark(t *testing.T) {
	_ = t
}

// Valid: assigning to a different field is fine.
func BenchmarkOtherField(b *testing.B) {
	_ = b.N
	for i := 0; i < b.N; i++ {
		// do work
	}
}
