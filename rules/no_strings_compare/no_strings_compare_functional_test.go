package no_strings_compare_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_strings_compare"
)

func TestNoStringsCompare(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_strings_compare", &no_strings_compare.NoStringsCompareRule{})
}
