package use_printf_suffix_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_printf_suffix"
)

func TestUsePrintfSuffix(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_printf_suffix", &use_printf_suffix.UsePrintfSuffixRule{})
}
