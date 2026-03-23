package no_naked_return_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_naked_return"
)

func TestNoNakedReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_naked_return", &no_naked_return.NoNakedReturnRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
