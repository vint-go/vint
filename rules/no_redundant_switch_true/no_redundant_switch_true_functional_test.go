package no_redundant_switch_true_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_switch_true"
)

func TestNoRedundantSwitchTrue(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_switch_true", &no_redundant_switch_true.NoRedundantSwitchTrueRule{})
}
