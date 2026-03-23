package no_impossible_condition_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_impossible_condition"
)

func TestNoImpossibleCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_condition", &no_impossible_condition.NoImpossibleConditionRule{})
}
