package no_variable_shadowing_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_variable_shadowing"
)

func TestNoVariableShadowing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_shadowing", &no_variable_shadowing.NoVariableShadowingRule{})
}
