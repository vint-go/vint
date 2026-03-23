package no_redundant_slice_expression_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_slice_expression"
)

func TestNoRedundantSliceExpression(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_slice_expression", &no_redundant_slice_expression.NoRedundantSliceExpressionRule{})
}
