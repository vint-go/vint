package no_single_iteration_loop_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_single_iteration_loop"
)

func TestNoSingleIterationLoop(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_iteration_loop", &no_single_iteration_loop.NoSingleIterationLoopRule{})
}
