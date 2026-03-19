package no_excessive_control_nesting_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_control_nesting"
)

func TestNoExcessiveControlNestingDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_control_nesting_default", &no_excessive_control_nesting.MaxControlNestingRule{}, &lint.RuleConfig{})
}

func TestNoExcessiveControlNesting(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_control_nesting", &no_excessive_control_nesting.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(2)},
	})
}
