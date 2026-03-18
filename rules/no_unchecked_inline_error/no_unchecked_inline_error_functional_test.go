package no_unchecked_inline_error_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unchecked_inline_error"
)

func TestNoUncheckedInlineError(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unchecked_inline_error", &no_unchecked_inline_error.NoUncheckedInlineErrorRule{})
}
