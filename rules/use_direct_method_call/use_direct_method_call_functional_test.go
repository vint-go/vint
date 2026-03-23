package use_direct_method_call_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_direct_method_call"
)

func TestUseDirectMethodCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_direct_method_call", &use_direct_method_call.UseDirectMethodCallRule{})
}
