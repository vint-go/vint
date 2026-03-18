package no_invalid_flag_name_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_invalid_flag_name"
)

func TestNoInvalidFlagName(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_flag_name", &no_invalid_flag_name.NoInvalidFlagNameRule{})
}
