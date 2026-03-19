package no_identical_branches_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_identical_branches"
)

func TestNoIdenticalBranches(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_identical_branches", &no_identical_branches.IdenticalBranchesRule{})
}
