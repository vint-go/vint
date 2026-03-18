package no_useless_math_call_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_useless_math_call"
)

func TestNoUselessMathCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_useless_math_call", &no_useless_math_call.NoUselessMathCallRule{})
}
