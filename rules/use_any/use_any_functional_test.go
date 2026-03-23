package use_any_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_any"
)

func TestUseAny(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_any", &use_any.UseAnyRule{})
}
