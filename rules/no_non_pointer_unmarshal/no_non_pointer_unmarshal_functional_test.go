package no_non_pointer_unmarshal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_non_pointer_unmarshal"
)

func TestNoNonPointerUnmarshal(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_pointer_unmarshal", &no_non_pointer_unmarshal.NoNonPointerUnmarshalRule{})
}
