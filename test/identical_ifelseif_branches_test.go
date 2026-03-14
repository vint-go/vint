package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestIdenticalIfElseIfBranches(t *testing.T) {
	testRule(t, "identical_ifelseif_branches", &rule.IdenticalIfElseIfBranchesRule{})
}
