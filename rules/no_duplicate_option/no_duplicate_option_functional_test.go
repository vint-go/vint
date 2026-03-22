package no_duplicate_option_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_option"
)

func TestNoDuplicateOption(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_option", &no_duplicate_option.NoDuplicateOptionRule{})
}
