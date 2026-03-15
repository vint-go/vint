package no_malformed_test_function_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_malformed_test_function"
)

func TestNoMalformedTestFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_test_function", &no_malformed_test_function.NoMalformedTestFunctionRule{})
}
