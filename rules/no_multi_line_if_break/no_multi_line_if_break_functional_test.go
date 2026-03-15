package no_multi_line_if_break_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_multi_line_if_break"
)

func TestNoMultiLineIfBreak(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_multi_line_if_break", &no_multi_line_if_break.NoMultiLineIfBreakRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"multi-if": true}},
	})
}
