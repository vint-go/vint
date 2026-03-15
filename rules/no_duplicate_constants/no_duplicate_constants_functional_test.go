package no_duplicate_constants_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_constants"
)

func TestNoDuplicateConstants(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_constants", &no_duplicate_constants.NoDuplicateConstantsRule{})
}
