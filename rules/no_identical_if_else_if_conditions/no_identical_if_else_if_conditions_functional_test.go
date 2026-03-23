package no_identical_if_else_if_conditions_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_identical_if_else_if_conditions"
)

func TestNoIdenticalIfElseIfConditions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_identical_if_else_if_conditions", &no_identical_if_else_if_conditions.IdenticalIfElseIfConditionsRule{})
}
