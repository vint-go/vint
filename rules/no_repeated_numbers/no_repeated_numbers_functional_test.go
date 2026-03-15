package no_repeated_numbers_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_repeated_numbers"
)

func TestNoRepeatedNumbers(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_repeated_numbers", &no_repeated_numbers.NoRepeatedNumbersRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"min-occurrences": int64(3),
			},
		},
	})
}
