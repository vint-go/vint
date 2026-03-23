package no_infinite_recursion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_infinite_recursion"
)

func TestNoInfiniteRecursion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_infinite_recursion", &no_infinite_recursion.NoInfiniteRecursionRule{})
}
