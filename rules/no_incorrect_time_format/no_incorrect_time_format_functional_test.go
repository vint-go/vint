package no_incorrect_time_format_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_incorrect_time_format"
)

func TestNoIncorrectTimeFormat(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_incorrect_time_format", &no_incorrect_time_format.NoIncorrectTimeFormatRule{})
}
