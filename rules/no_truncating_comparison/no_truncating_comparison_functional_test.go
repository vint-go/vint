package no_truncating_comparison_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_truncating_comparison"
)

func TestNoTruncatingComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_truncating_comparison", &no_truncating_comparison.NoTruncatingComparisonRule{})
}
