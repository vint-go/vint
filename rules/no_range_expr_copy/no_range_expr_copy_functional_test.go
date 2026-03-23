package no_range_expr_copy_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_range_expr_copy"
)

func TestNoRangeExprCopy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_expr_copy", &no_range_expr_copy.NoRangeExprCopyRule{})
}
