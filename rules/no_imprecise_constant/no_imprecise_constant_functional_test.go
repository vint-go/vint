package no_imprecise_constant_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_imprecise_constant"
)

func TestNoImpreciseConstant(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_imprecise_constant", &no_imprecise_constant.NoImpreciseConstantRule{})
}
