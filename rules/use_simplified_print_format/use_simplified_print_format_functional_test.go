package use_simplified_print_format_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_simplified_print_format"
)

func TestUseSimplifiedPrintFormat(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_print_format", &use_simplified_print_format.UseSimplifiedPrintFormatRule{})
}
