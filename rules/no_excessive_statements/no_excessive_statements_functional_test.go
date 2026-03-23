package no_excessive_statements_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_statements"
)

func TestNoExcessiveStatements(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_statements", &no_excessive_statements.NoExcessiveStatementsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
