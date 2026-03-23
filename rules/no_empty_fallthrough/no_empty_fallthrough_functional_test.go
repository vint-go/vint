package no_empty_fallthrough_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_empty_fallthrough"
)

func TestNoEmptyFallthrough(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_fallthrough", &no_empty_fallthrough.NoEmptyFallthroughRule{})
}
