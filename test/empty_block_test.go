package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestEmptyBlock(t *testing.T) {
	testRule(t, "empty_block", &rule.EmptyBlockRule{})
}
