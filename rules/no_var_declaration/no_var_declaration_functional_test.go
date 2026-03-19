package no_var_declaration_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_var_declaration"
)

func TestNoVarDeclaration(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_var_declaration", &no_var_declaration.VarDeclarationsRule{})
}
