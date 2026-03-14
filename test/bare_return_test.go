package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestBareReturn(t *testing.T) {
	testRule(t, "bare_return", &rule.BareReturnRule{})
}
