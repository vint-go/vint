package no_unused_parameter_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_parameter"
)

func TestNoUnusedParameter(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_parameter", &no_unused_parameter.NoUnusedParameterRule{})
}
