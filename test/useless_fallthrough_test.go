package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUselessFallthrough(t *testing.T) {
	testRule(t, "useless_fallthrough", &rule.UselessFallthroughRule{})
}
