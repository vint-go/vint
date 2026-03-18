package no_test_main_without_exit_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_test_main_without_exit"
)

func TestNoTestMainWithoutExit(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_test_main_without_exit", &no_test_main_without_exit.NoTestMainWithoutExitRule{})
}
