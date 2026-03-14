package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestGetReturn(t *testing.T) {
	testRule(t, "get_return", &rule.GetReturnRule{})
}
