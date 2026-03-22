package no_printf_format_mismatch_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_printf_format_mismatch"
)

func TestNoPrintfFormatMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_printf_format_mismatch", &no_printf_format_mismatch.NoPrintfFormatMismatchRule{})
}

func TestNoPrintfFormatMismatchCustomFuncs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_printf_format_mismatch_custom_funcs", &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"funcs": []any{"os.Setenv"},
		}},
	})
}
