package no_nan_comparison_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nan_comparison"
)

func TestNoNanComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nan_comparison", &no_nan_comparison.NoNanComparisonRule{})
}
