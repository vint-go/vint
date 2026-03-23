package use_slices_sort_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_slices_sort"
)

func TestUseSlicesSort(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slices_sort", &use_slices_sort.UseSlicesSort{})
}
