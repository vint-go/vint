package no_superfluous_else_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_superfluous_else"
)

func TestNoSuperfluousElse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_superfluous_else", &no_superfluous_else.SuperfluousElseRule{})
	functional_test_helpers.TestRule(t, "no_superfluous_else_scope", &no_superfluous_else.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	functional_test_helpers.TestRule(t, "no_superfluous_else_scope", &no_superfluous_else.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}
