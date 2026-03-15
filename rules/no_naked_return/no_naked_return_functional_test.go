package no_naked_return_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_naked_return"
)

func TestNoNakedReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_naked_return", &no_naked_return.NoNakedReturnRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
