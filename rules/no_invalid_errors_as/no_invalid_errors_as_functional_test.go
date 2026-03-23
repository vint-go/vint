package no_invalid_errors_as_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_errors_as"
)

func TestNoInvalidErrorsAs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_errors_as", &no_invalid_errors_as.NoInvalidErrorsAsRule{})
}
