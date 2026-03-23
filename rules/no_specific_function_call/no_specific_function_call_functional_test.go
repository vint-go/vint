package no_specific_function_call_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_specific_function_call"
)

func TestNoSpecificFunctionCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_specific_function_call", &no_specific_function_call.NoSpecificFunctionCallRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"name": "println"}},
	})
}
