package no_redundant_rune_conversion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_rune_conversion"
)

func TestNoRedundantRuneConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_rune_conversion", &no_redundant_rune_conversion.NoRedundantRuneConversionRule{})
}
