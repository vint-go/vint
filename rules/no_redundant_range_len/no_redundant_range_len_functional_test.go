package no_redundant_range_len_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_range_len"
)

func TestNoRedundantRangeLen(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_range_len", &no_redundant_range_len.NoRedundantRangeLenRule{})
}
