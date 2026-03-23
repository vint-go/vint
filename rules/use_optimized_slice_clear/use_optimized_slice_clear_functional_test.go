package use_optimized_slice_clear_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_optimized_slice_clear"
)

func TestUseOptimizedSliceClear(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_optimized_slice_clear", &use_optimized_slice_clear.UseOptimizedSliceClearRule{})
}
