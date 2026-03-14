package no_deprecated_hash_function_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_deprecated_hash_function"
)

func TestNoDeprecatedHashFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deprecated_hash_function", &no_deprecated_hash_function.NoDeprecatedHashFunctionRule{})
}
