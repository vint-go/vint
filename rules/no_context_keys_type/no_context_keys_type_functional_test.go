package no_context_keys_type_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_context_keys_type"
)

func TestNoContextKeysType(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_context_keys_type", &no_context_keys_type.ContextKeysType{})
}
