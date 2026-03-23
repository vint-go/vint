package no_nolint_non_machine_readable_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nolint_non_machine_readable"
)

func TestNoNolintNonMachineReadable(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nolint_non_machine_readable", &no_nolint_non_machine_readable.NoNolintNonMachineReadableRule{})
}
