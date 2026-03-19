package no_identical_switch_conditions_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_identical_switch_conditions"
)

func TestNoIdenticalSwitchConditions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_identical_switch_conditions", &no_identical_switch_conditions.IdenticalSwitchConditionsRule{})
}
