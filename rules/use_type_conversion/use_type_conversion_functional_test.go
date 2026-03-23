package use_type_conversion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_type_conversion"
)

func TestUseTypeConversion(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_conversion", &use_type_conversion.UseTypeConversionRule{})
}
