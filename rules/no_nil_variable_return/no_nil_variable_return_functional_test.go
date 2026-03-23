package no_nil_variable_return_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nil_variable_return"
)

func TestNoNilVariableReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_variable_return", &no_nil_variable_return.NoNilVariableReturnRule{})
}
