package no_useless_break_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_useless_break"
)

func TestNoUselessBreak(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_useless_break", &no_useless_break.UselessBreak{})
}
