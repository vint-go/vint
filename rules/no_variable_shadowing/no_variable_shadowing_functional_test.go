package no_variable_shadowing_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_variable_shadowing"
)

func TestNoVariableShadowing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_shadowing", &no_variable_shadowing.NoVariableShadowingRule{})
}

func TestNoVariableShadowingNonStrict(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_shadowing_non_strict", &no_variable_shadowing.NoVariableShadowingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"strict": false}},
	})
}
