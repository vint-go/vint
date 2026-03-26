package no_slice_bounds_out_of_range_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_slice_bounds_out_of_range"
)

func TestNoSliceBoundsOutOfRange(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_slice_bounds_out_of_range", &no_slice_bounds_out_of_range.NoSliceBoundsOutOfRangeRule{})
}
