package use_string_conversion_in_print_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_string_conversion_in_print"
)

func TestUseStringConversionInPrint(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_conversion_in_print", &use_string_conversion_in_print.UseStringConversionInPrintRule{})
}
