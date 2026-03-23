package no_non_wrapping_errorf_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_non_wrapping_errorf"
)

func TestNoNonWrappingErrorf(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_wrapping_errorf", &no_non_wrapping_errorf.NoNonWrappingErrorfRule{})
}
