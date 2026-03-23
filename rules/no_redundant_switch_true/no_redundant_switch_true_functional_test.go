package no_redundant_switch_true_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_switch_true"
)

func TestNoRedundantSwitchTrue(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_switch_true", &no_redundant_switch_true.NoRedundantSwitchTrueRule{})
}
