package no_dynamic_errors_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_dynamic_errors"
)

func TestNoDynamicErrors(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dynamic_errors", &no_dynamic_errors.NoDynamicErrorsRule{})
}
