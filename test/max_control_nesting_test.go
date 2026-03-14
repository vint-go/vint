package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestMaxControlNestingDefault(t *testing.T) {
	testRule(t, "max_control_nesting_default", &rule.MaxControlNestingRule{}, &lint.RuleConfig{})
}

func TestMaxControlNesting(t *testing.T) {
	testRule(t, "max_control_nesting", &rule.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(2)},
	})
}
