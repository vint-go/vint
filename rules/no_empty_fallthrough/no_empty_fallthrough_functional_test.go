package no_empty_fallthrough_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_empty_fallthrough"
)

func TestNoEmptyFallthrough(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_fallthrough", &no_empty_fallthrough.NoEmptyFallthroughRule{})
}
