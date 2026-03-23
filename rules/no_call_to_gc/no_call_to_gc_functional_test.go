package no_call_to_gc_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_call_to_gc"
)

func TestNoCallToGC(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_call_to_gc", &no_call_to_gc.CallToGCRule{})
}
