package no_excessive_shift_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_shift"
)

func TestNoExcessiveShift(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_shift", &no_excessive_shift.NoExcessiveShiftRule{})
}
