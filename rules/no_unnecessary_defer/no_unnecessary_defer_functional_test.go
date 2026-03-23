package no_unnecessary_defer_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unnecessary_defer"
)

func TestNoUnnecessaryDefer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_defer", &no_unnecessary_defer.NoUnnecessaryDeferRule{})
}
