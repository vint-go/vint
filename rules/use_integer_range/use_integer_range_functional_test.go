package use_integer_range_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_integer_range"
)

func TestUseIntegerRange(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_integer_range", &use_integer_range.UseIntegerRangeRule{})
}
