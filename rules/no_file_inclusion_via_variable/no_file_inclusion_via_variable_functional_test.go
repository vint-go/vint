package no_file_inclusion_via_variable_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_file_inclusion_via_variable"
)

func TestNoFileInclusionViaVariable(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_file_inclusion_via_variable", &no_file_inclusion_via_variable.NoFileInclusionViaVariableRule{})
}
