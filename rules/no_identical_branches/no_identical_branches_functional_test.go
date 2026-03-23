package no_identical_branches_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_identical_branches"
)

func TestNoIdenticalBranches(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_identical_branches", &no_identical_branches.IdenticalBranchesRule{})
}
