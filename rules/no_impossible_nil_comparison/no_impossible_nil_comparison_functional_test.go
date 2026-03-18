package no_impossible_nil_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_impossible_nil_comparison"
)

func TestNoImpossibleNilComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_nil_comparison", &no_impossible_nil_comparison.NoImpossibleNilComparisonRule{})
}
