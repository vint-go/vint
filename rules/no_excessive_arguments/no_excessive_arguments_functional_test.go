package no_excessive_arguments_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_arguments"
)

func TestNoExcessiveArgumentsDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_arguments_default", &no_excessive_arguments.NoExcessiveArgumentsRule{})
}

func TestNoExcessiveArguments(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_arguments", &no_excessive_arguments.NoExcessiveArgumentsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}
