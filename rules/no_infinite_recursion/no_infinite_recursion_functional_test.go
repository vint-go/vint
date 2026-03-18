package no_infinite_recursion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_infinite_recursion"
)

func TestNoInfiniteRecursion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_infinite_recursion", &no_infinite_recursion.NoInfiniteRecursionRule{})
}
