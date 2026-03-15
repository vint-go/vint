package no_constant_parameter_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_constant_parameter"
)

func TestNoConstantParameter(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_constant_parameter", &no_constant_parameter.NoConstantParameterRule{})
}
