package no_redundant_range_val_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_range_val"
)

func TestNoRedundantRangeVal(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_range_val", &no_redundant_range_val.RangeRule{})
}
