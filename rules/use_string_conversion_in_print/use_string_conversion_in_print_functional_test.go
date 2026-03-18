package use_string_conversion_in_print_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_string_conversion_in_print"
)

func TestUseStringConversionInPrint(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_conversion_in_print", &use_string_conversion_in_print.UseStringConversionInPrintRule{})
}
