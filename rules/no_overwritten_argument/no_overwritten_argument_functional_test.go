package no_overwritten_argument_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_overwritten_argument"
)

func TestNoOverwrittenArgument(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_overwritten_argument", &no_overwritten_argument.NoOverwrittenArgumentRule{})
}
