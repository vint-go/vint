package use_errors_new_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_errors_new"
)

func TestUseErrorsNew(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_errors_new", &use_errors_new.UseErrorsNewRule{})
}
