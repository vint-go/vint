package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestIdenticalSwitchBranches(t *testing.T) {
	testRule(t, "identical_switch_branches", &rule.IdenticalSwitchBranchesRule{})
}
