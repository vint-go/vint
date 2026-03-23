package no_reflect_value_compare_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_reflect_value_compare"
)

func TestNoReflectValueCompare(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_reflect_value_compare", &no_reflect_value_compare.NoReflectValueCompareRule{})
}
