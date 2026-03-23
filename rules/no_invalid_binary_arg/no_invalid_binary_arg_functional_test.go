package no_invalid_binary_arg_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_binary_arg"
)

func TestNoInvalidBinaryArg(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_binary_arg", &no_invalid_binary_arg.NoInvalidBinaryArgRule{})
}
