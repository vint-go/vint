package no_non_pointer_unmarshal_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_non_pointer_unmarshal"
)

func TestNoNonPointerUnmarshal(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_pointer_unmarshal", &no_non_pointer_unmarshal.NoNonPointerUnmarshalRule{})
}
