package no_address_of_dereference_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_address_of_dereference"
)

func TestNoAddressOfDereference(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_address_of_dereference", &no_address_of_dereference.NoAddressOfDereferenceRule{})
}
