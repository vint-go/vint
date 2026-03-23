package test_test

import (
	"testing"

	"github.com/vint-go/vint/rule"
)

func TestEmptyBlock(t *testing.T) {
	testRule(t, "empty_block", &rule.EmptyBlockRule{})
}
