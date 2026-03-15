package no_single_arg_append_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_single_arg_append"
)

func TestNoSingleArgAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_arg_append", &no_single_arg_append.NoSingleArgAppendRule{})
}
