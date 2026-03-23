package no_single_case_switch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_single_case_switch"
)

func TestNoSingleCaseSwitch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_case_switch", &no_single_case_switch.NoSingleCaseSwitchRule{})
}
