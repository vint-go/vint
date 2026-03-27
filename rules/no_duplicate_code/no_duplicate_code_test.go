package no_duplicate_code_test

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/no_duplicate_code"
)

// helper creates lint.File objects from testdata filenames.
func loadFile(t *testing.T, name string) *lint.File {
	t.Helper()
	path := filepath.Join("testdata", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	fset := token.NewFileSet()
	pkg := lint.NewPackage(fset, lint.DefaultGoVersion, lint.NewSafeSharedImporter())
	file, err := pkg.AddFile(path, content)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", path, err)
	}
	return file
}

func TestCrossFileDuplicateDetected(t *testing.T) {
	rule := &no_duplicate_code.NoDuplicateCodeRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	fileA := loadFile(t, "cross_dup_a.go")
	fileB := loadFile(t, "cross_dup_b.go")

	rule.Collect(fileA, nil)
	rule.Collect(fileB, nil)

	failures := rule.Finalize()
	if len(failures) == 0 {
		t.Fatal("expected at least one duplicate code failure across files, got none")
	}

	// Verify that failures reference both files (one in position, one in message).
	foundCrossFile := false
	for _, f := range failures {
		t.Logf("failure: %s (at %s)", f.Failure, f.Position.Start.Filename)
		if f.Category != lint.FailureCategoryComplexity {
			t.Errorf("expected category %q, got %q", lint.FailureCategoryComplexity, f.Category)
		}
		// The position should reference one file, the message should reference the other.
		posFile := f.Position.Start.Filename
		msgRefA := contains(f.Failure, "cross_dup_a.go")
		msgRefB := contains(f.Failure, "cross_dup_b.go")
		if (contains(posFile, "cross_dup_b.go") && msgRefA) ||
			(contains(posFile, "cross_dup_a.go") && msgRefB) {
			foundCrossFile = true
		}
		// Verify that line numbers are set.
		if f.Position.Start.Line == 0 {
			t.Error("expected start line to be set")
		}
		if f.Position.End.Line == 0 {
			t.Error("expected end line to be set")
		}
	}
	if !foundCrossFile {
		t.Error("expected at least one failure referencing both cross_dup_a.go and cross_dup_b.go")
	}
}

func TestNoDuplicateNonMatching(t *testing.T) {
	rule := &no_duplicate_code.NoDuplicateCodeRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	fileA := loadFile(t, "no_dup_a.go")
	fileB := loadFile(t, "no_dup_b.go")

	rule.Collect(fileA, nil)
	rule.Collect(fileB, nil)

	failures := rule.Finalize()
	if len(failures) != 0 {
		for _, f := range failures {
			t.Logf("unexpected failure: %s", f.Failure)
		}
		t.Fatalf("expected no failures for non-matching files, got %d", len(failures))
	}
}

func TestSingleFileNoDuplicate(t *testing.T) {
	rule := &no_duplicate_code.NoDuplicateCodeRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	file := loadFile(t, "cross_dup_a.go")
	rule.Collect(file, nil)

	failures := rule.Finalize()
	if len(failures) != 0 {
		t.Fatalf("expected no failures for single file, got %d", len(failures))
	}
}

func containsBoth(s, a, b string) bool {
	return contains(s, a) && contains(s, b)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
