package no_excessive_statements_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_statements"
)

func TestNoExcessiveStatements(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_statements", &no_excessive_statements.NoExcessiveStatementsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
