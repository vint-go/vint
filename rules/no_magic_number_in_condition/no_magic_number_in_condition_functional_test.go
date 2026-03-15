package no_magic_number_in_condition_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_magic_number_in_condition"
)

func TestNoMagicNumberInCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_condition", &no_magic_number_in_condition.NoMagicNumberInConditionRule{})
}
