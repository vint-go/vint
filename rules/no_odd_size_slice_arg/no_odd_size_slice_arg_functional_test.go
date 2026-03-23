package no_odd_size_slice_arg_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_odd_size_slice_arg"
)

func TestNoOddSizeSliceArg(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_odd_size_slice_arg", &no_odd_size_slice_arg.NoOddSizeSliceArgRule{})
}
