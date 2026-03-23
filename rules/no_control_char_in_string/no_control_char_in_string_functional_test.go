package no_control_char_in_string_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_control_char_in_string"
)

func TestNoControlCharInString(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_control_char_in_string", &no_control_char_in_string.NoControlCharInStringRule{})
}
