package no_immediate_new_deref_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_immediate_new_deref"
)

func TestNoImmediateNewDeref(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_immediate_new_deref", &no_immediate_new_deref.NoImmediateNewDerefRule{})
}
