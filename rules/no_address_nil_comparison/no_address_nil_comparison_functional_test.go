package no_address_nil_comparison_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_address_nil_comparison"
)

func TestNoAddressNilComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_address_nil_comparison", &no_address_nil_comparison.NoAddressNilComparisonRule{})
}
