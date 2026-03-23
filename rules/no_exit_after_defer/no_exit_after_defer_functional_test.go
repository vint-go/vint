package no_exit_after_defer_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_exit_after_defer"
)

func TestNoExitAfterDefer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_exit_after_defer", &no_exit_after_defer.NoExitAfterDeferRule{})
}
