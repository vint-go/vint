package no_explicit_bool_comparison_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_explicit_bool_comparison"
)

func TestNoExplicitBoolComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_explicit_bool_comparison", &no_explicit_bool_comparison.NoExplicitBoolComparisonRule{})
}
