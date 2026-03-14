package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUseAny(t *testing.T) {
	testRule(t, "use_any", &rule.UseAnyRule{})
}
