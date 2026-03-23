package no_integer_division_truncation_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_integer_division_truncation"
)

func TestNoIntegerDivisionTruncation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_integer_division_truncation", &no_integer_division_truncation.NoIntegerDivisionTruncationRule{})
}
