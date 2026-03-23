package no_unused_type_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_type"
)

func TestNoUnusedType(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_type", &no_unused_type.NoUnusedTypeRule{})
}
