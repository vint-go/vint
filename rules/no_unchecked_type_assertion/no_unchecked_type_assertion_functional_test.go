package no_unchecked_type_assertion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unchecked_type_assertion"
)

func TestNoUncheckedTypeAssertion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unchecked_type_assertion", &no_unchecked_type_assertion.NoUncheckedTypeAssertionRule{})
}
