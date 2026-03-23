package no_unnecessary_type_parens_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unnecessary_type_parens"
)

func TestNoUnnecessaryTypeParens(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_type_parens", &no_unnecessary_type_parens.NoUnnecessaryTypeParensRule{})
}
