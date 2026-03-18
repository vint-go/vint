package no_nil_context_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nil_context"
)

func TestNoNilContext(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_context", &no_nil_context.NoNilContextRule{})
}
