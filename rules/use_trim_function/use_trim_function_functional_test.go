package use_trim_function_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_trim_function"
)

func TestUseTrimFunction(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_trim_function", &use_trim_function.UseTrimFunctionRule{})
}
