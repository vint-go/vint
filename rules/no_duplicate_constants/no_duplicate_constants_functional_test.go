package no_duplicate_constants_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_constants"
)

func TestNoDuplicateConstants(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_constants", &no_duplicate_constants.NoDuplicateConstantsRule{})
}
