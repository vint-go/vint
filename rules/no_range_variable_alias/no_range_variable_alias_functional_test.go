package no_range_variable_alias_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_range_variable_alias"
)

func TestNoRangeVariableAlias(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_variable_alias", &no_range_variable_alias.NoRangeVariableAliasRule{})
}
