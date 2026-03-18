package no_invariant_loop_condition_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_invariant_loop_condition"
)

func TestNoInvariantLoopCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invariant_loop_condition", &no_invariant_loop_condition.NoInvariantLoopConditionRule{})
}
