package no_multi_line_func_break_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_multi_line_func_break"
)

func TestNoMultiLineFuncBreak(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_multi_line_func_break", &no_multi_line_func_break.NoMultiLineFuncBreakRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"multi-func": true}},
	})
}
