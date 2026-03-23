package no_excessive_function_results_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_function_results"
)

func TestNoExcessiveFunctionResults(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_function_results", &no_excessive_function_results.FunctionResultsLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}
