package no_duplicate_code

import (
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/strowk/vint/lint"
)

func benchFile(b *testing.B, filename string) *lint.File {
	b.Helper()
	src, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		b.Fatal(err)
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		b.Fatal(err)
	}
	return &lint.File{
		Name: filename,
		AST:  astFile,
	}
}

// BenchmarkNoDuplicateCode_Matching benchmarks the worst case for duplicate
// detection: 30 structurally-identical large functions (all pairs match).
func BenchmarkNoDuplicateCode_Matching(b *testing.B) {
	file := benchFile(b, "bench_large.go")

	rule := &NoDuplicateCodeRule{}
	_ = rule.Configure(nil)

	failures := rule.Apply(file, nil)
	if len(failures) == 0 {
		b.Fatal("expected failures from bench_large.go but got none")
	}
	b.Logf("rule found %d failures", len(failures))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Apply(file, nil)
	}
}

// BenchmarkNoDuplicateCode_NonMatching benchmarks the real-world bottleneck:
// 15 large but structurally-different functions. No pairs match the threshold,
// so the DP runs to full completion for every pair (105 comparisons).
func BenchmarkNoDuplicateCode_NonMatching(b *testing.B) {
	file := benchFile(b, "bench_non_matching.go")

	rule := &NoDuplicateCodeRule{}
	_ = rule.Configure(nil)

	failures := rule.Apply(file, nil)
	if len(failures) != 0 {
		b.Fatalf("expected 0 failures from bench_non_matching.go but got %d", len(failures))
	}
	b.Log("rule found 0 failures (expected)")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Apply(file, nil)
	}
}
