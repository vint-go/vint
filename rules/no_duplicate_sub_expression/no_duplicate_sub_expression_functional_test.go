package no_duplicate_sub_expression_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_sub_expression"
)

func TestNoDuplicateSubExpression(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_sub_expression", &no_duplicate_sub_expression.NoDuplicateSubExpressionRule{})
}
