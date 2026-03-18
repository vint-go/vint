package no_empty_branch_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_empty_branch"
)

func TestNoEmptyBranch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_branch", &no_empty_branch.NoEmptyBranchRule{})
}
