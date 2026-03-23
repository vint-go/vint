package no_excessive_arguments_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_arguments"
)

func TestNoExcessiveArgumentsDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_arguments_default", &no_excessive_arguments.NoExcessiveArgumentsRule{})
}

func TestNoExcessiveArguments(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_arguments", &no_excessive_arguments.NoExcessiveArgumentsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}
