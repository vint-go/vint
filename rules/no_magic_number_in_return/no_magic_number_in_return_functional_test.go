package no_magic_number_in_return_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_magic_number_in_return"
)

func TestNoMagicNumberInReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_return", &no_magic_number_in_return.NoMagicNumberInReturnRule{})
}
