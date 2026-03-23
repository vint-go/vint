package use_error_last_return_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_error_last_return"
)

func TestUseErrorLastReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_error_last_return", &use_error_last_return.UseErrorLastReturnRule{})
}
