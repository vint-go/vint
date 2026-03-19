package no_unconditional_recursion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unconditional_recursion"
)

func TestNoUnconditionalRecursion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unconditional_recursion", &no_unconditional_recursion.UnconditionalRecursionRule{})
}
