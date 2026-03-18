package no_suspicious_sort_slice_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_suspicious_sort_slice"
)

func TestNoSuspiciousSortSlice(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_suspicious_sort_slice", &no_suspicious_sort_slice.NoSuspiciousSortSliceRule{})
}
