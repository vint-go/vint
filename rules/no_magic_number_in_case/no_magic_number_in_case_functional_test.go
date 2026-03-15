package no_magic_number_in_case_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_magic_number_in_case"
)

func TestNoMagicNumberInCase(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_case", &no_magic_number_in_case.NoMagicNumberInCaseRule{})
}
