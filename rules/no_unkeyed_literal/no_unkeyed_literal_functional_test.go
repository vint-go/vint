package no_unkeyed_literal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unkeyed_literal"
)

func TestNoUnkeyedLiteral(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unkeyed_literal", &no_unkeyed_literal.NoUnkeyedLiteralRule{})
}
