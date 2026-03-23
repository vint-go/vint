package no_excessive_shift_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_shift"
)

func TestNoExcessiveShift(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_shift", &no_excessive_shift.NoExcessiveShiftRule{})
}
