package no_non_wrapping_errorf_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_non_wrapping_errorf"
)

func TestNoNonWrappingErrorf(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_wrapping_errorf", &no_non_wrapping_errorf.NoNonWrappingErrorfRule{})
}
