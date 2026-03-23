package no_trailing_blank_line_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_trailing_blank_line"
)

func TestNoTrailingBlankLine(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_trailing_blank_line", &no_trailing_blank_line.NoTrailingBlankLineRule{})
}
