package no_invalid_utf8_string_arg_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_utf8_string_arg"
)

func TestNoInvalidUtf8StringArg(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_utf8_string_arg", &no_invalid_utf8_string_arg.NoInvalidUtf8StringArgRule{})
}
