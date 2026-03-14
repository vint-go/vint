package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestIdenticalSwitchConditions(t *testing.T) {
	testRule(t, "identical_switch_conditions", &rule.IdenticalSwitchConditionsRule{})
}
