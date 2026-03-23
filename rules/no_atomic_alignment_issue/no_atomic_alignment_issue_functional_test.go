package no_atomic_alignment_issue_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_atomic_alignment_issue"
)

func TestNoAtomicAlignmentIssue(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_atomic_alignment_issue", &no_atomic_alignment_issue.NoAtomicAlignmentIssueRule{})
}
