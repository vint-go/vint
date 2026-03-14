package no_direct_error_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_direct_error_comparison"
)

func TestNoDirectErrorComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_direct_error_comparison", &no_direct_error_comparison.NoDirectErrorComparisonRule{})
}
