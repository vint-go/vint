package use_direct_return_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_direct_return"
)

func TestUseDirectReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_direct_return", &use_direct_return.IfReturnRule{})
}
