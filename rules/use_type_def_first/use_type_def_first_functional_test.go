package use_type_def_first_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_type_def_first"
)

func TestUseTypeDefFirst(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_def_first", &use_type_def_first.UseTypeDefFirstRule{})
}
