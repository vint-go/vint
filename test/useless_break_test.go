package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUselessBreak(t *testing.T) {
	testRule(t, "useless_break", &rule.UselessBreak{})
}
