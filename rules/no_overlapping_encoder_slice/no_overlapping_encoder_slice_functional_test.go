package no_overlapping_encoder_slice_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_overlapping_encoder_slice"
)

func TestNoOverlappingEncoderSlice(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_overlapping_encoder_slice", &no_overlapping_encoder_slice.NoOverlappingEncoderSliceRule{})
}
