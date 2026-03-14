package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestIdenticalIfElseIfConditions(t *testing.T) {
	testRule(t, "identical_ifelseif_conditions", &rule.IdenticalIfElseIfConditionsRule{})
}
