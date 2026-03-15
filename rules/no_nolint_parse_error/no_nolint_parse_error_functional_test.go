package no_nolint_parse_error_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nolint_parse_error"
)

func TestNoNolintParseError(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nolint_parse_error", &no_nolint_parse_error.NoNolintParseErrorRule{})
}
