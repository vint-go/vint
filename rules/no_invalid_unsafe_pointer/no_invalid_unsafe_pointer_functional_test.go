package no_invalid_unsafe_pointer_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_invalid_unsafe_pointer"
)

func TestNoInvalidUnsafePointer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_unsafe_pointer", &no_invalid_unsafe_pointer.NoInvalidUnsafePointerRule{})
}
