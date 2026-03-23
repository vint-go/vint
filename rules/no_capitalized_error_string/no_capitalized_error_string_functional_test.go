package no_capitalized_error_string_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_capitalized_error_string"
)

func TestNoCapitalizedErrorString(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_capitalized_error_string", &no_capitalized_error_string.NoCapitalizedErrorStringRule{})
}
