package no_explicit_bool_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_explicit_bool_comparison"
)

func TestNoExplicitBoolComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_explicit_bool_comparison", &no_explicit_bool_comparison.NoExplicitBoolComparisonRule{})
}
