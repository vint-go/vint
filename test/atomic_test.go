package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestAtomic(t *testing.T) {
	testRule(t, "atomic", &rule.AtomicRule{})
}
