package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

// Test that left and right side of Binary operators (only AND, OR) are swapable.
func TestOptimizeOperandsOrder(t *testing.T) {
	testRule(t, "optimize_operands_order", &rule.OptimizeOperandsOrderRule{}, &lint.RuleConfig{})
}
