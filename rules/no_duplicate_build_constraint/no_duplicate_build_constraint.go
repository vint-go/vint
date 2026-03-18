package no_duplicate_build_constraint

import (
	"fmt"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateBuildConstraintRule detects duplicate build constraints in the same file.
// Having multiple identical //go:build lines is redundant and likely a copy-paste mistake.
type NoDuplicateBuildConstraintRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateBuildConstraintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect all //go:build constraint texts and track which we've seen.
	seen := map[string]bool{}

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			if !strings.HasPrefix(text, "//go:build") {
				continue
			}

			// Normalize: trim the prefix and whitespace to get the constraint expression.
			expr := strings.TrimSpace(strings.TrimPrefix(text, "//go:build"))
			if expr == "" {
				continue
			}

			if seen[expr] {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("duplicate build constraint: %q", text),
				})
			} else {
				seen[expr] = true
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoDuplicateBuildConstraintRule) Name() string {
	return "noDuplicateBuildConstraint"
}

// Group returns the rule group.
func (*NoDuplicateBuildConstraintRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateBuildConstraintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
