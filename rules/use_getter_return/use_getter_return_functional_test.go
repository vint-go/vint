package use_getter_return_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_getter_return"
)

func TestUseGetterReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_getter_return", &use_getter_return.GetReturnRule{})
}
