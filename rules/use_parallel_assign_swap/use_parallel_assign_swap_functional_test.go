package use_parallel_assign_swap_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_parallel_assign_swap"
)

func TestUseParallelAssignSwap(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_parallel_assign_swap", &use_parallel_assign_swap.UseParallelAssignSwapRule{})
}
