package no_defer_in_infinite_loop_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_defer_in_infinite_loop"
)

func TestNoDeferInInfiniteLoop(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_in_infinite_loop", &no_defer_in_infinite_loop.NoDeferInInfiniteLoopRule{})
}
