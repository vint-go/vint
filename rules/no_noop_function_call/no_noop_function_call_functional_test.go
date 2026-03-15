package no_noop_function_call_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_noop_function_call"
)

func TestNoNoopFunctionCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_noop_function_call", &no_noop_function_call.NoNoopFunctionCallRule{})
}
