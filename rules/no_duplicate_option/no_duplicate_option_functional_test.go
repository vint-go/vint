package no_duplicate_option_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_option"
)

func TestNoDuplicateOption(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_option", &no_duplicate_option.NoDuplicateOptionRule{})
}
