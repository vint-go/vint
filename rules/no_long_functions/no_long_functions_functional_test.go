package no_long_functions_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_long_functions"
)

func TestNoLongFunctions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_long_functions", &no_long_functions.NoLongFunctionsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
