package no_nolint_without_specific_linter_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nolint_without_specific_linter"
)

func TestNoNolintWithoutSpecificLinter(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nolint_without_specific_linter", &no_nolint_without_specific_linter.NoNolintWithoutSpecificLinterRule{})
}
