package no_excessive_function_results_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_function_results"
)

func TestNoExcessiveFunctionResults(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_function_results", &no_excessive_function_results.FunctionResultsLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}
