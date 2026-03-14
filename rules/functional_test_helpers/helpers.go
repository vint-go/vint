// Package functional_test_helpers provides shared test utilities for
// packaged rule functional tests.
package functional_test_helpers

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/strowk/vint/lint"
)

// TestRule is a functional test helper that:
//  1. Reads a testdata Go file
//  2. Parses MATCH instructions from comments (format: // MATCH /expected failure message/)
//  3. Runs the linter with the given rule
//  4. Verifies that failures match expectations
func TestRule(tb testing.TB, filename string, rule lint.Rule, ruleConfig ...*lint.RuleConfig) {
	tb.Helper()

	filePath := filepath.Join("testdata", filename+".go")
	src, err := os.ReadFile(filePath)
	if err != nil {
		tb.Fatalf("Bad filename path in test for %s: %v", rule.Name(), err)
	}

	// Configure the rule if it implements ConfigurableRule
	if cr, ok := rule.(lint.ConfigurableRule); ok {
		var args lint.Arguments
		if len(ruleConfig) > 0 && ruleConfig[0] != nil {
			args = ruleConfig[0].Arguments
		}
		if err := cr.Configure(args); err != nil {
			tb.Fatalf("Cannot configure rule %s: %v", rule.Name(), err)
		}
	}

	ins := parseInstructions(tb, filePath, src)
	if ins == nil {
		assertSuccess(tb, filePath, rule)
		return
	}
	assertFailures(tb, filePath, rule, ins)
}

func assertSuccess(tb testing.TB, filePath string, rule lint.Rule) {
	tb.Helper()

	l := lint.New(os.ReadFile, 0)
	ps, err := l.Lint([][]string{{filePath}}, []lint.Rule{rule}, lint.Config{})
	if err != nil {
		tb.Errorf("Linting %s: %v", filePath, err)
		return
	}

	failures := ""
	for p := range ps {
		failures += p.Failure
	}
	if failures != "" {
		tb.Errorf("Expected the rule to pass but got the following failures: %s", failures)
	}
}

func assertFailures(tb testing.TB, filePath string, rule lint.Rule, ins []instruction) {
	tb.Helper()

	l := lint.New(os.ReadFile, 0)
	ps, err := l.Lint([][]string{{filePath}}, []lint.Rule{rule}, lint.Config{})
	if err != nil {
		tb.Errorf("Linting %s: %v", filePath, err)
		return
	}

	failures := []lint.Failure{}
	for f := range ps {
		failures = append(failures, f)
	}

	type simplifiedFailure struct {
		File    string
		Line    int
		Failure string
	}
	var reportedFailures []simplifiedFailure

	for _, in := range ins {
		ok := false
		for i, p := range failures {
			if p.Position.Start.Line != in.Line {
				continue
			}
			if in.Match == p.Failure {
				copy(failures[i:], failures[i+1:])
				failures = failures[:len(failures)-1]
				ok = true
				break
			}
		}
		if !ok {
			reportedFailures = append(reportedFailures, simplifiedFailure{
				File:    filePath,
				Line:    in.Line,
				Failure: "want: " + in.Match,
			})
		}
	}

	for _, p := range failures {
		reportedFailures = append(reportedFailures, simplifiedFailure{
			File:    filePath,
			Line:    p.Position.Start.Line,
			Failure: "got:  " + p.Failure,
		})
	}

	slices.SortFunc(reportedFailures, func(a, b simplifiedFailure) int {
		if a.File != b.File {
			return strings.Compare(a.File, b.File)
		}
		if a.Line != b.Line {
			return a.Line - b.Line
		}
		return strings.Compare(a.Failure, b.Failure)
	})

	errorMessage := ""
	lastFileLine := ""
	for _, p := range reportedFailures {
		currentFileLine := fmt.Sprintf("%s:%d", p.File, p.Line)
		if currentFileLine != lastFileLine {
			if errorMessage != "" {
				tb.Error(errorMessage)
			}
			errorMessage = fmt.Sprintf("problem at %s: ", currentFileLine)
			lastFileLine = currentFileLine
		}
		errorMessage += "\n" + p.Failure
	}

	if errorMessage != "" {
		tb.Error(errorMessage)
	}
}

type instruction struct {
	Line  int
	Match string
}

// parseInstructions parses MATCH instructions from comments in a Go source file.
// Returns nil if none were found.
func parseInstructions(tb testing.TB, filename string, src []byte) []instruction {
	tb.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		tb.Fatalf("Test file %v does not parse: %v", filename, err)
	}

	var ins []instruction
	for _, cg := range f.Comments {
		ln := fset.Position(cg.Pos()).Line
		raw := cg.Text()
		for line := range strings.SplitSeq(raw, "\n") {
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "ignore") {
				continue
			}
			if line == "OK" && ins == nil {
				ins = []instruction{}
				continue
			}
			if !strings.Contains(line, "MATCH") {
				continue
			}
			match, err := extractPattern(line)
			if err != nil {
				tb.Fatalf("At %v:%d: %v", filename, ln, err)
			}
			ins = append(ins, instruction{
				Line:  ln,
				Match: match,
			})
		}
	}
	return ins
}

func extractPattern(line string) (string, error) {
	a, b := strings.Index(line, "/"), strings.LastIndex(line, "/")
	if a == -1 || a == b {
		return "", fmt.Errorf("malformed match instruction %q", line)
	}
	return line[a+1 : b], nil
}
