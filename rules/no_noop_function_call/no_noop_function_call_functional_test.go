package no_noop_function_call_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_noop_function_call"
)

func TestNoNoopFunctionCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_noop_function_call", &no_noop_function_call.NoNoopFunctionCallRule{})
}
