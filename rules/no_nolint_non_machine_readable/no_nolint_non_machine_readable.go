package no_nolint_non_machine_readable

import (
	"fmt"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNolintNonMachineReadableRule reports //nolint directives that are not written
// in a machine-readable format (i.e. have a space between // and nolint).
type NoNolintNonMachineReadableRule struct{}

// Apply applies the rule to given file.
func (r *NoNolintNonMachineReadableRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// We are looking for comments like "// nolint" (with a space).
			// The machine-readable format is "//nolint" (no space).
			// We need to detect "// nolint" but not "//nolint".
			if !strings.HasPrefix(text, "// nolint") {
				continue
			}

			// Extract the rest after "// " to build the suggested fix
			rest := strings.TrimPrefix(text, "// ")

			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       comment,
				Failure:    fmt.Sprintf("nolint directive %q is not machine-readable, use %q", text, "//"+rest),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoNolintNonMachineReadableRule) Name() string {
	return "noNolintNonMachineReadable"
}

// Group returns the rule group.
func (*NoNolintNonMachineReadableRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoNolintNonMachineReadableRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
