package no_unsafe_deserialization_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unsafe_deserialization"
)

func TestNoUnsafeDeserialization(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsafe_deserialization", &no_unsafe_deserialization.NoUnsafeDeserializationRule{})
}
