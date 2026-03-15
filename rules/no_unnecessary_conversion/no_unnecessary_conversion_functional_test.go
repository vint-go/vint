package no_unnecessary_conversion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_conversion"
)

func TestNoUnnecessaryConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_conversion", &no_unnecessary_conversion.NoUnnecessaryConversionRule{})
}
