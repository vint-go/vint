package no_redundant_test_main_exit_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_test_main_exit"
)

func TestNoRedundantTestMainExit(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_test_main_exit_test", &no_redundant_test_main_exit.RedundantTestMainExitRule{})
}
