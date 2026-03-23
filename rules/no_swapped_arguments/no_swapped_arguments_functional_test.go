package no_swapped_arguments_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_swapped_arguments"
)

func TestNoSwappedArguments(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_swapped_arguments", &no_swapped_arguments.NoSwappedArgumentsRule{})
}
