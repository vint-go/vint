package no_yoda_condition_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_yoda_condition"
)

func TestNoYodaCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_yoda_condition", &no_yoda_condition.NoYodaConditionRule{})
}
