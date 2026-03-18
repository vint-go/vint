package no_unnecessary_defer_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_defer"
)

func TestNoUnnecessaryDefer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_defer", &no_unnecessary_defer.NoUnnecessaryDeferRule{})
}
