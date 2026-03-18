package use_idiomatic_error_name_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_idiomatic_error_name"
)

func TestUseIdiomaticErrorName(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_error_name", &use_idiomatic_error_name.UseIdiomaticErrorNameRule{})
}
