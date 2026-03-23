package use_error_naming_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_error_naming"
)

func TestUseErrorNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_error_naming", &use_error_naming.ErrorNamingRule{})
}
