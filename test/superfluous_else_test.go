package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestSuperfluousElse(t *testing.T) {
	testRule(t, "superfluous_else", &rule.SuperfluousElseRule{})
	testRule(t, "superfluous_else_scope", &rule.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "superfluous_else_scope", &rule.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}
