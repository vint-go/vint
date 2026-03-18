package no_inappropriate_context_key_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_inappropriate_context_key"
)

func TestNoInappropriateContextKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_inappropriate_context_key", &no_inappropriate_context_key.NoInappropriateContextKeyRule{})
}
