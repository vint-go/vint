package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestIfReturn(t *testing.T) {
	testRule(t, "if_return", &rule.IfReturnRule{})
}
