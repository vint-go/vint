package no_unkeyed_literal_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unkeyed_literal"
)

func TestNoUnkeyedLiteral(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unkeyed_literal", &no_unkeyed_literal.NoUnkeyedLiteralRule{})
}
