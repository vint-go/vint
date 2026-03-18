package use_time_equal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_time_equal"
)

func TestUseTimeEqual(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_equal", &use_time_equal.UseTimeEqualRule{})
}
