package no_single_arg_append_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_single_arg_append"
)

func TestNoSingleArgAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_arg_append", &no_single_arg_append.NoSingleArgAppendRule{})
}
