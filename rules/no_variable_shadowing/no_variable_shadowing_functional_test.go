package no_variable_shadowing_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_variable_shadowing"
)

func TestNoVariableShadowing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_shadowing", &no_variable_shadowing.NoVariableShadowingRule{})
}

func TestNoVariableShadowingNonStrict(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_shadowing_non_strict", &no_variable_shadowing.NoVariableShadowingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"strict": false}},
	})
}
