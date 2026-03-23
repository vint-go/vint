package no_deprecated_hash_function_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_deprecated_hash_function"
)

func TestNoDeprecatedHashFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deprecated_hash_function", &no_deprecated_hash_function.NoDeprecatedHashFunctionRule{})
}
