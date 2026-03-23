package no_inappropriate_context_key_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_inappropriate_context_key"
)

func TestNoInappropriateContextKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_inappropriate_context_key", &no_inappropriate_context_key.NoInappropriateContextKeyRule{})
}
