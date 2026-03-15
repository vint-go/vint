package no_nil_func_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nil_func_comparison"
)

func TestNoNilFuncComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_func_comparison", &no_nil_func_comparison.NoNilFuncComparisonRule{})
}
