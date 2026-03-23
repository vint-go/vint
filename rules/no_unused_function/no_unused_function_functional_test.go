package no_unused_function_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_function"
)

func TestNoUnusedFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_function", &no_unused_function.NoUnusedFunctionRule{})
}
