package no_pointer_to_ref_param_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_pointer_to_ref_param"
)

func TestNoPointerToRefParam(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_pointer_to_ref_param", &no_pointer_to_ref_param.NoPointerToRefParamRule{})
}
