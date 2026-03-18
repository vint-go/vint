package no_dubious_bit_shift_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_dubious_bit_shift"
)

func TestNoDubiousBitShift(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dubious_bit_shift", &no_dubious_bit_shift.NoDubiousBitShiftRule{})
}
