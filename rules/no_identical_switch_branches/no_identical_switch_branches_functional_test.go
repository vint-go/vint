package no_identical_switch_branches_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_identical_switch_branches"
)

func TestNoIdenticalSwitchBranches(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_identical_switch_branches", &no_identical_switch_branches.IdenticalSwitchBranchesRule{})
}
