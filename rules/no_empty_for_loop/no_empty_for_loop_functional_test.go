package no_empty_for_loop_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_empty_for_loop"
)

func TestNoEmptyForLoop(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_for_loop", &no_empty_for_loop.NoEmptyForLoopRule{})
}
