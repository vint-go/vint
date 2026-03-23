package no_implicit_const_value_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_implicit_const_value"
)

func TestNoImplicitConstValue(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_implicit_const_value", &no_implicit_const_value.NoImplicitConstValueRule{})
}
