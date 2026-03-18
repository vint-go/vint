package use_bytes_equal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_bytes_equal"
)

func TestUseBytesEqual(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_bytes_equal", &use_bytes_equal.UseBytesEqualRule{})
}
