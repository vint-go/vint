package use_errors_as_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_errors_as"
)

func TestUseErrorsAs(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_errors_as", &use_errors_as.UseErrorsAsRule{})
}
