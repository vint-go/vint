package no_duplicate_argument_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_argument"
)

func TestNoDuplicateArgument(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_argument", &no_duplicate_argument.NoDuplicateArgumentRule{})
}
