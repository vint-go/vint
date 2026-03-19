package no_call_to_gc_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_call_to_gc"
)

func TestNoCallToGC(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_call_to_gc", &no_call_to_gc.CallToGCRule{})
}
