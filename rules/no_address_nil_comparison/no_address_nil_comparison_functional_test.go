package no_address_nil_comparison_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_address_nil_comparison"
)

func TestNoAddressNilComparison(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_address_nil_comparison", &no_address_nil_comparison.NoAddressNilComparisonRule{})
}
