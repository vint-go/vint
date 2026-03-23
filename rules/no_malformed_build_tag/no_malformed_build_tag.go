package no_malformed_build_tag

import (
	"fmt"
	"go/build/constraint"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMalformedBuildTagRule checks build tags for correctness.
// It verifies that //go:build and // +build directives are well-formed,
// consistent, and placed before the package clause.
type NoMalformedBuildTagRule struct{}

// Apply applies the rule to given file.
func (r *NoMalformedBuildTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	packageLine := 0
	if file.AST.Package.IsValid() {
		packageLine = file.ToPosition(file.AST.Package).Line
	}

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// Check //go:build directives
			if strings.HasPrefix(text, "//go:build") {
				commentLine := file.ToPosition(comment.Pos()).Line

				// Check if build tag is after the package clause
				if packageLine > 0 && commentLine > packageLine {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryBadPractice,
						Confidence: 1,
						Node:       comment,
						Failure:    "build tag must appear before the package clause",
					})
					continue
				}

				// Check if the constraint expression is valid
				expr := strings.TrimPrefix(text, "//go:build")
				expr = strings.TrimSpace(expr)
				if expr == "" {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryBadPractice,
						Confidence: 1,
						Node:       comment,
						Failure:    "build tag has empty constraint expression",
					})
					continue
				}

				_, err := constraint.Parse(text)
				if err != nil {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryBadPractice,
						Confidence: 1,
						Node:       comment,
						Failure:    fmt.Sprintf("malformed build tag: %s", err),
					})
				}
			}

			// Check // +build directives
			if strings.HasPrefix(text, "// +build") || text == "//+build" || strings.HasPrefix(text, "//+build ") {
				commentLine := file.ToPosition(comment.Pos()).Line

				// Check if build tag is after the package clause
				if packageLine > 0 && commentLine > packageLine {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryBadPractice,
						Confidence: 1,
						Node:       comment,
						Failure:    "build tag must appear before the package clause",
					})
				}
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoMalformedBuildTagRule) Name() string {
	return "noMalformedBuildTag"
}

// Group returns the rule group.
func (*NoMalformedBuildTagRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMalformedBuildTagRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
