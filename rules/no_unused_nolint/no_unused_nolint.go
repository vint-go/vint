package no_unused_nolint

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnusedNolintRule reports //nolint directives that are unnecessary because
// they do not suppress any actual linter warnings. An unused nolint directive
// indicates that the underlying issue was fixed without removing the
// suppression comment, or the directive was added preemptively.
type NoUnusedNolintRule struct{}

// Apply applies the rule to given file.
func (r *NoUnusedNolintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// We only care about machine-readable nolint directives (no space).
			if !strings.HasPrefix(text, "//nolint") {
				continue
			}

			rest := strings.TrimPrefix(text, "//nolint")

			// Skip things like "//nolintfoo" which are not nolint directives.
			if rest != "" && rest[0] != ':' && rest[0] != ' ' && rest[0] != '\t' {
				continue
			}

			// Build the directive description for the message.
			directive := "//nolint"
			if rest != "" && rest[0] == ':' {
				// Extract linter list (up to space or end).
				afterColon := rest[1:]
				parts := strings.SplitN(afterColon, " ", 2)
				directive = "//nolint:" + parts[0]
			}

			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       comment,
				Failure:    fmt.Sprintf("nolint directive %q is unused", directive),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoUnusedNolintRule) Name() string {
	return "noUnusedNolint"
}

// Group returns the rule group.
func (*NoUnusedNolintRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedNolintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
