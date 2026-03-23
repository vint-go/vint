package no_redundant_var_type_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_var_type"
)

func TestNoRedundantVarType(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_var_type", &no_redundant_var_type.NoRedundantVarTypeRule{})
}
