package no_duplicate_if_condition_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_if_condition"
)

func TestNoDuplicateIfCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_if_condition", &no_duplicate_if_condition.NoDuplicateIfConditionRule{})
}
