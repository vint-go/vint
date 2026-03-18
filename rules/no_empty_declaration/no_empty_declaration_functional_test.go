package no_empty_declaration_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_empty_declaration"
)

func TestNoEmptyDeclaration(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_declaration", &no_empty_declaration.NoEmptyDeclarationRule{})
}
