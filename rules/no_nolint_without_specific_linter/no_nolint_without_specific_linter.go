package no_nolint_without_specific_linter

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNolintWithoutSpecificLinterRule reports //nolint directives that do not
// specify which linter(s) are being suppressed. A bare //nolint suppresses all
// linters for the annotated line, which is overly broad and can mask real issues.
type NoNolintWithoutSpecificLinterRule struct{}

// Apply applies the rule to given file.
func (r *NoNolintWithoutSpecificLinterRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// We only care about machine-readable nolint directives (no space).
			// The "// nolint" variant (with space) is handled by
			// noNolintNonMachineReadable, so skip those here.
			if !strings.HasPrefix(text, "//nolint") {
				continue
			}

			// After "//nolint", the next character tells us what we have:
			// - ':' means linters are specified (e.g. //nolint:errcheck)
			// - ' ' or end of string means bare //nolint
			// - any other letter means it's a different directive (e.g. //nolintfoo)
			rest := strings.TrimPrefix(text, "//nolint")

			if rest == "" {
				// Bare "//nolint" with nothing after
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("nolint directive %q should mention specific linter(s), e.g. //nolint:errcheck", text),
				})
				continue
			}

			if rest[0] == ':' {
				// Has a colon, so linters are specified. This is fine.
				continue
			}

			if rest[0] == ' ' || rest[0] == '\t' {
				// Bare "//nolint" followed by a comment (e.g. "//nolint // reason")
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("nolint directive %q should mention specific linter(s), e.g. //nolint:errcheck", text),
				})
				continue
			}

			// Otherwise it's something like "//nolintfoo" which is not a
			// nolint directive at all; ignore it.
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoNolintWithoutSpecificLinterRule) Name() string {
	return "noNolintWithoutSpecificLinter"
}

// Group returns the rule group.
func (*NoNolintWithoutSpecificLinterRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoNolintWithoutSpecificLinterRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
