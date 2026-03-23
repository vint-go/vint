package no_misplaced_default_case_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_misplaced_default_case"
)

func TestNoMisplacedDefaultCase(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_misplaced_default_case", &no_misplaced_default_case.NoMisplacedDefaultCaseRule{})
}
