package no_truncating_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_truncating_comparison"
)

func TestNoTruncatingComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_truncating_comparison", &no_truncating_comparison.NoTruncatingComparisonRule{})
}
