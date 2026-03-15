package no_deep_equal_errors_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_deep_equal_errors"
)

func TestNoDeepEqualErrors(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deep_equal_errors", &no_deep_equal_errors.NoDeepEqualErrorsRule{})
}
