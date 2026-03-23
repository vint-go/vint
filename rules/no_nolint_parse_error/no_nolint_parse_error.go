package no_nolint_parse_error

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// fullDirectivePattern validates the full nolint directive.
// Valid forms:
//
//	//nolint
//	//nolint:linter1,linter2
//	//nolint:linter1,linter2 // explanation
//	//nolint // explanation
//
// Linter names consist of word characters and hyphens: [\w-]+
var fullDirectivePattern = regexp.MustCompile(`^//nolint(?::(\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*))?\s*(//.*)?$`)

// NoNolintParseErrorRule reports //nolint directives that have malformed syntax
// and cannot be properly parsed.
type NoNolintParseErrorRule struct{}

// Apply applies the rule to given file.
func (r *NoNolintParseErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// Only consider machine-readable nolint directives (no space).
			// The "// nolint" variant is handled by noNolintNonMachineReadable.
			if !strings.HasPrefix(text, "//nolint") {
				continue
			}

			rest := strings.TrimPrefix(text, "//nolint")

			// Not a nolint directive at all (e.g. "//nolintfoo").
			if rest != "" && rest[0] != ':' && rest[0] != ' ' && rest[0] != '\t' {
				continue
			}

			// Bare "//nolint" or "//nolint // explanation" are valid
			// (handled by noNolintWithoutSpecificLinter for enforcement).
			// We only care about parse errors here.

			// Try to match the full directive pattern.
			if fullDirectivePattern.MatchString(text) {
				continue
			}

			// The directive did not match the valid pattern => parse error.
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       comment,
				Failure:    fmt.Sprintf("malformed nolint directive %q: expected format is //nolint[:<linter1>,<linter2>] [// explanation]", text),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoNolintParseErrorRule) Name() string {
	return "noNolintParseError"
}

// Group returns the rule group.
func (*NoNolintParseErrorRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNolintParseErrorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
