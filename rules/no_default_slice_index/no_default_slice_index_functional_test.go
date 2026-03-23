package no_default_slice_index_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_default_slice_index"
)

func TestNoDefaultSliceIndex(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_default_slice_index", &no_default_slice_index.NoDefaultSliceIndexRule{})
}
