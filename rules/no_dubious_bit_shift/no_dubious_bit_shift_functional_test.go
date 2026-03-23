package no_dubious_bit_shift_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_dubious_bit_shift"
)

func TestNoDubiousBitShift(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dubious_bit_shift", &no_dubious_bit_shift.NoDubiousBitShiftRule{})
}
