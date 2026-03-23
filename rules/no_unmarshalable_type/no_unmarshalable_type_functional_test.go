package no_unmarshalable_type_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unmarshalable_type"
)

func TestNoUnmarshalableType(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unmarshalable_type", &no_unmarshalable_type.NoUnmarshalableTypeRule{})
}
