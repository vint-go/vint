package no_integer_overflow_conversion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_integer_overflow_conversion"
)

func TestNoIntegerOverflowConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_integer_overflow_conversion", &no_integer_overflow_conversion.NoIntegerOverflowConversionRule{})
}
