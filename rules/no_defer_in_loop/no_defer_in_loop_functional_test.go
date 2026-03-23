package no_defer_in_loop_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_defer_in_loop"
)

func TestNoDeferInLoop(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_in_loop", &no_defer_in_loop.NoDeferInLoopRule{})
}
