package no_constant_parameter_test

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/no_constant_parameter"
)

func loadFile(t *testing.T, name string) *lint.File {
	t.Helper()
	path := filepath.Join("testdata", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	fset := token.NewFileSet()
	pkg := lint.NewPackage(fset, lint.DefaultGoVersion, lint.NewSharedImporter())
	file, err := pkg.AddFile(path, content)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", path, err)
	}
	return file
}

func TestNoConstantParameter(t *testing.T) {
	rule := &no_constant_parameter.NoConstantParameterRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	file := loadFile(t, "no_constant_parameter.go")
	rule.Collect(file, nil)
	failures := rule.Finalize()

	expected := map[int]string{
		9:  "times always receives 3",
		24: "verbose always receives true",
		39: `separator always receives ","`,
	}

	assertExpectedFailures(t, failures, expected, "no_constant_parameter.go")
}

// TestCrossFileConstantParameter verifies that call sites in other files
// prevent false positives. When file1 calls process("fast") and file2 calls
// process("slow"), the parameter should NOT be flagged.
func TestCrossFileConstantParameter(t *testing.T) {
	rule := &no_constant_parameter.NoConstantParameterRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	fileA := loadFile(t, "cross_file_a.go")
	fileB := loadFile(t, "cross_file_b.go")

	rule.Collect(fileA, nil)
	rule.Collect(fileB, nil)

	failures := rule.Finalize()
	if len(failures) != 0 {
		for _, f := range failures {
			t.Logf("unexpected failure at line %d: %s", f.Position.Start.Line, f.Failure)
		}
		t.Fatalf("expected no failures for cross-file varying parameters, got %d", len(failures))
	}
}

// TestCrossFileConstantParameterStillFlags verifies that when ALL call sites
// across files pass the same constant, the parameter is still flagged.
func TestCrossFileConstantParameterStillFlags(t *testing.T) {
	rule := &no_constant_parameter.NoConstantParameterRule{}
	if err := rule.Configure(nil); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	fileA := loadFile(t, "cross_file_same_a.go")
	fileB := loadFile(t, "cross_file_same_b.go")

	rule.Collect(fileA, nil)
	rule.Collect(fileB, nil)

	failures := rule.Finalize()

	expected := map[int]string{
		5: `mode always receives "fast"`,
	}

	assertExpectedFailures(t, failures, expected, "cross_file_same_a.go")
}

func assertExpectedFailures(t *testing.T, failures []lint.Failure, expected map[int]string, filename string) {
	t.Helper()

	matched := map[int]bool{}
	var unexpected []lint.Failure

	for _, f := range failures {
		line := f.Position.Start.Line
		if expMsg, ok := expected[line]; ok {
			if f.Failure != expMsg {
				t.Errorf("line %d: expected %q, got %q", line, expMsg, f.Failure)
			}
			matched[line] = true
		} else {
			unexpected = append(unexpected, f)
		}
	}

	for line, msg := range expected {
		if !matched[line] {
			// Check if the failure is in the expected file but at a different location
			found := false
			for _, f := range failures {
				if strings.Contains(f.Position.Start.Filename, filename) && f.Failure == msg {
					t.Errorf("expected failure at line %d but found it at line %d: %s", line, f.Position.Start.Line, msg)
					found = true
					break
				}
			}
			if !found {
				t.Errorf("missing expected failure at line %d: %s", line, msg)
			}
		}
	}

	for _, f := range unexpected {
		t.Errorf("unexpected failure at line %d: %s", f.Position.Start.Line, f.Failure)
	}
}
