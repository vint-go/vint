package no_useless_fallthrough_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_useless_fallthrough"
)

func TestNoUselessFallthrough(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_useless_fallthrough", &no_useless_fallthrough.UselessFallthroughRule{})
}
