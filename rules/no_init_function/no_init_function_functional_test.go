package no_init_function_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_init_function"
)

func TestNoInitFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_init_function", &no_init_function.NoInitFunctionRule{})
}
