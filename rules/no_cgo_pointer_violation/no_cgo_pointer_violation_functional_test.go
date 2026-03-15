package no_cgo_pointer_violation_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_cgo_pointer_violation"
)

func TestNoCgoPointerViolation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_cgo_pointer_violation", &no_cgo_pointer_violation.NoCgoPointerViolationRule{})
}
