package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestCallToGC(t *testing.T) {
	testRule(t, "call_to_gc", &rule.CallToGCRule{})
}
