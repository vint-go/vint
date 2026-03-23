package no_flag_parameter_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_flag_parameter"
)

func TestNoFlagParameter(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_flag_parameter", &no_flag_parameter.FlagParamRule{})
}
