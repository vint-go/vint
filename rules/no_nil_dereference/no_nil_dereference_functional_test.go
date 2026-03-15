package no_nil_dereference_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nil_dereference"
)

func TestNoNilDereference(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_dereference", &no_nil_dereference.NoNilDereferenceRule{})
}
