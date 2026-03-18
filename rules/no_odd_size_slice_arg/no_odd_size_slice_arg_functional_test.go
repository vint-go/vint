package no_odd_size_slice_arg_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_odd_size_slice_arg"
)

func TestNoOddSizeSliceArg(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_odd_size_slice_arg", &no_odd_size_slice_arg.NoOddSizeSliceArgRule{})
}
