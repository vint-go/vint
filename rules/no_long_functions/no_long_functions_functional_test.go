package no_long_functions_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_long_functions"
)

func TestNoLongFunctions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_long_functions", &no_long_functions.NoLongFunctionsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
