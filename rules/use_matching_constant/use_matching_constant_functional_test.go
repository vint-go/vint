package use_matching_constant_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_matching_constant"
)

func TestUseMatchingConstant(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_matching_constant", &use_matching_constant.UseMatchingConstantRule{})
}
