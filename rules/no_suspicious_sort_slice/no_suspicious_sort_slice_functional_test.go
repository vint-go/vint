package no_suspicious_sort_slice_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_suspicious_sort_slice"
)

func TestNoSuspiciousSortSlice(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_suspicious_sort_slice", &no_suspicious_sort_slice.NoSuspiciousSortSliceRule{})
}
