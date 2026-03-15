package no_magic_number_in_argument_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_magic_number_in_argument"
)

func TestNoMagicNumberInArgument(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_argument", &no_magic_number_in_argument.NoMagicNumberInArgumentRule{})
}
