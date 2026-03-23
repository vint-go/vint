package no_ineffective_break_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ineffective_break"
)

func TestNoIneffectiveBreak(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ineffective_break", &no_ineffective_break.NoIneffectiveBreakRule{})
}
