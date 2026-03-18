package no_unsigned_negative_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unsigned_negative_comparison"
)

func TestNoUnsignedNegativeComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsigned_negative_comparison", &no_unsigned_negative_comparison.NoUnsignedNegativeComparisonRule{})
}
