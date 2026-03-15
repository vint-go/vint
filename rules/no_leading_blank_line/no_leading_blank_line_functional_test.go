package no_leading_blank_line_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_leading_blank_line"
)

func TestNoLeadingBlankLine(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_leading_blank_line", &no_leading_blank_line.NoLeadingBlankLineRule{})
}
