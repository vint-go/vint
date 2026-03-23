package no_context_keys_type_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_context_keys_type"
)

func TestNoContextKeysType(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_context_keys_type", &no_context_keys_type.ContextKeysType{})
}
