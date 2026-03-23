package no_string_int_conversion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_string_int_conversion"
)

func TestNoStringIntConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_string_int_conversion", &no_string_int_conversion.NoStringIntConversionRule{})
}
