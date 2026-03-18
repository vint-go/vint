package no_unnecessary_type_parens_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_type_parens"
)

func TestNoUnnecessaryTypeParens(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_type_parens", &no_unnecessary_type_parens.NoUnnecessaryTypeParensRule{})
}
