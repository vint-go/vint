package no_impossible_condition_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_impossible_condition"
)

func TestNoImpossibleCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_condition", &no_impossible_condition.NoImpossibleConditionRule{})
}
