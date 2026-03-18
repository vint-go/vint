package no_mixed_case_hex_literal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_mixed_case_hex_literal"
)

func TestNoMixedCaseHexLiteral(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_mixed_case_hex_literal", &no_mixed_case_hex_literal.NoMixedCaseHexLiteralRule{})
}
