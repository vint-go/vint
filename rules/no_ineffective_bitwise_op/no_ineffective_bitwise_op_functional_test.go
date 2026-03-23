package no_ineffective_bitwise_op_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ineffective_bitwise_op"
)

func TestNoIneffectiveBitwiseOp(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ineffective_bitwise_op", &no_ineffective_bitwise_op.NoIneffectiveBitwiseOpRule{})
}
