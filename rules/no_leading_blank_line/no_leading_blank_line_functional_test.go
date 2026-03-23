package no_leading_blank_line_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_leading_blank_line"
)

func TestNoLeadingBlankLine(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_leading_blank_line", &no_leading_blank_line.NoLeadingBlankLineRule{})
}
