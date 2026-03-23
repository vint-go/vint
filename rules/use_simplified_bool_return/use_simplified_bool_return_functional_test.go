package use_simplified_bool_return_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_simplified_bool_return"
)

func TestUseSimplifiedBoolReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_bool_return", &use_simplified_bool_return.UseSimplifiedBoolReturnRule{})
}
