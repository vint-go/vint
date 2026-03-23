package use_error_method_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_error_method"
)

func TestUseErrorMethod(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_error_method", &use_error_method.UseErrorMethodRule{})
}
