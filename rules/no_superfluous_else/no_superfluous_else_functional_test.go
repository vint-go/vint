package no_superfluous_else_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_superfluous_else"
)

func TestNoSuperfluousElse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_superfluous_else", &no_superfluous_else.SuperfluousElseRule{})
	functional_test_helpers.TestRule(t, "no_superfluous_else_scope", &no_superfluous_else.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	functional_test_helpers.TestRule(t, "no_superfluous_else_scope", &no_superfluous_else.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}
