package no_nolint_without_explanation

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNolintWithoutExplanationRule reports //nolint directives that do not include
// an explanation for why the lint suppression is necessary.
type NoNolintWithoutExplanationRule struct {
	excludeLinters map[string]bool
}

// Configure validates and applies the rule configuration.
func (r *NoNolintWithoutExplanationRule) Configure(arguments lint.Arguments) error {
	r.excludeLinters = map[string]bool{}

	for _, arg := range arguments {
		cfg, ok := arg.(map[string]any)
		if !ok {
			continue
		}
		if v, ok := cfg["exclude"]; ok {
			if list, ok := v.([]any); ok {
				for _, item := range list {
					if s, ok := item.(string); ok {
						r.excludeLinters[s] = true
					}
				}
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoNolintWithoutExplanationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// We only care about machine-readable nolint directives (no space).
			if !strings.HasPrefix(text, "//nolint") {
				continue
			}

			rest := strings.TrimPrefix(text, "//nolint")

			// Must have a colon to specify linters; bare //nolint is handled
			// by noNolintWithoutSpecificLinter, so skip those here.
			if rest == "" || rest[0] != ':' {
				continue
			}

			// Extract the linter list (everything between ':' and the next ' ' or end)
			afterColon := rest[1:]

			// Split on " //" to separate linter list from the explanation comment
			parts := strings.SplitN(afterColon, " //", 2)
			linterPart := parts[0]

			// Check if all specified linters are in the exclude list
			if r.allLintersExcluded(linterPart) {
				continue
			}

			// Check for an explanation after the linter list
			if len(parts) < 2 {
				// No explanation at all
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("nolint directive %q should provide an explanation, e.g. //nolint:linter // reason", text),
				})
				continue
			}

			// There is a " //" part — check if the explanation is non-empty
			explanation := strings.TrimSpace(parts[1])
			if explanation == "" {
				// Bare trailing "//" with no text
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("nolint directive %q should provide an explanation, e.g. //nolint:linter // reason", text),
				})
			}
		}
	}

	return failures
}

// allLintersExcluded returns true if every linter mentioned in the comma-separated
// linter part is in the exclude set.
func (r *NoNolintWithoutExplanationRule) allLintersExcluded(linterPart string) bool {
	if len(r.excludeLinters) == 0 {
		return false
	}

	// The linter part may contain spaces after the linters
	// (e.g., "errcheck " if followed by " //"), but we already split on " //".
	// It may also contain multiple linters separated by commas.
	linters := strings.Split(linterPart, ",")
	for _, l := range linters {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if !r.excludeLinters[l] {
			return false
		}
	}
	return true
}

// Name returns the rule name.
func (*NoNolintWithoutExplanationRule) Name() string {
	return "noNolintWithoutExplanation"
}

// Group returns the rule group.
func (*NoNolintWithoutExplanationRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoNolintWithoutExplanationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
