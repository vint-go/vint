package no_unused_variable_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unused_variable"
)

func TestNoUnusedVariable(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_variable", &no_unused_variable.NoUnusedVariableRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"post-statements-are-reads": false,
				"local-variables-are-used":  true,
				"generated-is-used":         true,
			},
		},
	})
}
