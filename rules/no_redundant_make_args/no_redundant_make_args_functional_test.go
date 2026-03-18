package no_redundant_make_args_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_make_args"
)

func TestNoRedundantMakeArgs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_make_args", &no_redundant_make_args.NoRedundantMakeArgsRule{})
}
