package no_modified_parameter_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_modified_parameter"
)

func TestNoModifiedParameter(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_modified_parameter", &no_modified_parameter.ModifiesParamRule{})
}
