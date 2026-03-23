package no_unused_field_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_field"
)

func TestNoUnusedField(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_field", &no_unused_field.NoUnusedFieldRule{})
}
