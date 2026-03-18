package use_time_method_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_time_method"
)

func TestUseTimeMethod(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_method", &use_time_method.UseTimeMethodRule{})
}
