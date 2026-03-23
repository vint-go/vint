package no_redundant_string_concat_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_string_concat"
)

func TestNoRedundantStringConcat(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_string_concat", &no_redundant_string_concat.NoRedundantStringConcatRule{})
}
