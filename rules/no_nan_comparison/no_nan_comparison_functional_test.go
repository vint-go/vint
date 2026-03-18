package no_nan_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nan_comparison"
)

func TestNoNanComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nan_comparison", &no_nan_comparison.NoNanComparisonRule{})
}
