package no_redundant_make_args_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_make_args"
)

func TestNoRedundantMakeArgs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_make_args", &no_redundant_make_args.NoRedundantMakeArgsRule{})
}
