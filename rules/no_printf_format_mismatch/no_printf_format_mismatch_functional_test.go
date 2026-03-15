package no_printf_format_mismatch_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_printf_format_mismatch"
)

func TestNoPrintfFormatMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_printf_format_mismatch", &no_printf_format_mismatch.NoPrintfFormatMismatchRule{})
}
