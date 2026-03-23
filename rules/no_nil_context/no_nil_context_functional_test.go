package no_nil_context_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nil_context"
)

func TestNoNilContext(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_context", &no_nil_context.NoNilContextRule{})
}
