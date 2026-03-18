package no_implicit_const_value_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_implicit_const_value"
)

func TestNoImplicitConstValue(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_implicit_const_value", &no_implicit_const_value.NoImplicitConstValueRule{})
}
