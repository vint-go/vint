package no_zero_bytes_repeat_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_zero_bytes_repeat"
)

func TestNoZeroBytesRepeat(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_zero_bytes_repeat", &no_zero_bytes_repeat.NoZeroBytesRepeatRule{})
}
