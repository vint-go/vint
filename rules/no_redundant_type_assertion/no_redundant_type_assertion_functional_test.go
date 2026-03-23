package no_redundant_type_assertion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_type_assertion"
)

func TestNoRedundantTypeAssertion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_type_assertion", &no_redundant_type_assertion.NoRedundantTypeAssertionRule{})
}
