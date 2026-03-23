package use_direct_string_range_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_direct_string_range"
)

func TestUseDirectStringRange(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_direct_string_range", &use_direct_string_range.UseDirectStringRangeRule{})
}
