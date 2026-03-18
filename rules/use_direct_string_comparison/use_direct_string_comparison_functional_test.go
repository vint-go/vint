package use_direct_string_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_direct_string_comparison"
)

func TestUseDirectStringComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_direct_string_comparison", &use_direct_string_comparison.UseDirectStringComparisonRule{})
}
