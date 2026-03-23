package use_optimal_operands_order_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_optimal_operands_order"
)

func TestUseOptimalOperandsOrder(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_optimal_operands_order", &use_optimal_operands_order.OptimizeOperandsOrderRule{})
}
