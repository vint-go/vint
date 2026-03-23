package no_excessive_control_nesting_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_control_nesting"
)

func TestNoExcessiveControlNestingDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_control_nesting_default", &no_excessive_control_nesting.MaxControlNestingRule{}, &lint.RuleConfig{})
}

func TestNoExcessiveControlNesting(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_control_nesting", &no_excessive_control_nesting.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(2)},
	})
}
