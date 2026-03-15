package no_invalid_errors_as_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_invalid_errors_as"
)

func TestNoInvalidErrorsAs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_errors_as", &no_invalid_errors_as.NoInvalidErrorsAsRule{})
}
