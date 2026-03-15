package no_unused_function_result_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unused_function_result"
)

func TestNoUnusedFunctionResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_function_result", &no_unused_function_result.NoUnusedFunctionResultRule{})
}
