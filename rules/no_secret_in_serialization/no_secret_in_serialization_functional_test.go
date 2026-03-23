package no_secret_in_serialization_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_secret_in_serialization"
)

func TestNoSecretInSerialization(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_secret_in_serialization", &no_secret_in_serialization.NoSecretInSerializationRule{})
}
