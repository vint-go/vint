package no_test_fatal_in_goroutine_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_test_fatal_in_goroutine"
)

func TestNoTestFatalInGoroutine(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_test_fatal_in_goroutine", &no_test_fatal_in_goroutine.NoTestFatalInGoroutineRule{})
}
