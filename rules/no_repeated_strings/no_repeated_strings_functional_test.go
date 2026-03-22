package no_repeated_strings_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_repeated_strings"
)

func TestNoRepeatedStrings(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_repeated_strings", &no_repeated_strings.NoRepeatedStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"min-occurrences": int64(3),
				"min-length":     int64(3),
				"ignore-tests":   true,
			},
		},
	})
}

func TestNoRepeatedStringsIgnoreCallsDisabled(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_repeated_strings_ignore_calls_disabled", &no_repeated_strings.NoRepeatedStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"min-occurrences": int64(3),
				"min-length":     int64(3),
				"ignore-tests":   true,
				"ignore-calls":   false,
			},
		},
	})
}

func TestNoRepeatedStringsEvalConstExpressions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_repeated_strings_eval_const_expressions", &no_repeated_strings.NoRepeatedStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"min-occurrences":        int64(3),
				"min-length":            int64(3),
				"ignore-tests":          true,
				"eval-const-expressions": true,
			},
		},
	})
}
