package no_imprecise_constant_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_imprecise_constant"
)

func TestNoImpreciseConstant(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_imprecise_constant", &no_imprecise_constant.NoImpreciseConstantRule{})
}
