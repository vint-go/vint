package no_unsafe_deserialization_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unsafe_deserialization"
)

func TestNoUnsafeDeserialization(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsafe_deserialization", &no_unsafe_deserialization.NoUnsafeDeserializationRule{})
}
