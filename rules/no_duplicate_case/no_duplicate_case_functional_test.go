package no_duplicate_case_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_case"
)

func TestNoDuplicateCase(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_case", &no_duplicate_case.NoDuplicateCaseRule{})
}
