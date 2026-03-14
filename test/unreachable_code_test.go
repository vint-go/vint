package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUnreachableCode(t *testing.T) {
	testRule(t, "unreachable_code", &rule.UnreachableCodeRule{})
}
