package no_deep_exit_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_deep_exit"
)

func TestNoDeepExit(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deep_exit", &no_deep_exit.DeepExitRule{})
	functional_test_helpers.TestRule(t, "no_deep_exit_test", &no_deep_exit.DeepExitRule{})
}
