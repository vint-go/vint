package no_defer_in_loop_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_defer_in_loop"
)

func TestNoDeferInLoop(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_in_loop", &no_defer_in_loop.NoDeferInLoopRule{})
}
