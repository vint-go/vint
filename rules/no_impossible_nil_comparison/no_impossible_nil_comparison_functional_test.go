package no_impossible_nil_comparison_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_impossible_nil_comparison"
)

func TestNoImpossibleNilComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_nil_comparison", &no_impossible_nil_comparison.NoImpossibleNilComparisonRule{})
}
