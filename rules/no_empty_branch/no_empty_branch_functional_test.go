package no_empty_branch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_empty_branch"
)

func TestNoEmptyBranch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_branch", &no_empty_branch.NoEmptyBranchRule{})
}
