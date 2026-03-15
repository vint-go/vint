package no_multi_line_func_break_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_multi_line_func_break"
)

func TestNoMultiLineFuncBreak(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_multi_line_func_break", &no_multi_line_func_break.NoMultiLineFuncBreakRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"multi-func": true}},
	})
}
