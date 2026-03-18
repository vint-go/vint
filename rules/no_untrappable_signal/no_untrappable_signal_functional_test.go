package no_untrappable_signal_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_untrappable_signal"
)

func TestNoUntrappableSignal(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_untrappable_signal", &no_untrappable_signal.NoUntrappableSignalRule{})
}
