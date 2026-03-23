package no_line_too_long_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_line_too_long"
)

func TestNoLineTooLong(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_line_too_long", &no_line_too_long.NoLineTooLongRule{})
}
