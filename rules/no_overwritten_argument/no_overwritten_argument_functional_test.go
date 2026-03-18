package no_overwritten_argument_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_overwritten_argument"
)

func TestNoOverwrittenArgument(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_overwritten_argument", &no_overwritten_argument.NoOverwrittenArgumentRule{})
}
