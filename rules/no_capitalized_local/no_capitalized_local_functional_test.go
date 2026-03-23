package no_capitalized_local_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_capitalized_local"
)

func TestNoCapitalizedLocal(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_capitalized_local", &no_capitalized_local.NoCapitalizedLocalRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"paramsOnly": true}},
	})
}
