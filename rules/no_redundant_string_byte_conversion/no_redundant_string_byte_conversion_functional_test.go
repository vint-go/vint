package no_redundant_string_byte_conversion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_string_byte_conversion"
)

func TestNoRedundantStringByteConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_string_byte_conversion", &no_redundant_string_byte_conversion.NoRedundantStringByteConversionRule{})
}
