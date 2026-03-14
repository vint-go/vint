package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUnconditionalRecursion(t *testing.T) {
	testRule(t, "unconditional_recursion", &rule.UnconditionalRecursionRule{})
}
