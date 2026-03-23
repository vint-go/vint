package no_repeated_numbers_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_repeated_numbers"
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
