package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestRedefinesBuiltinID(t *testing.T) {
	testRule(t, "redefines_builtin_id", &rule.RedefinesBuiltinIDRule{})
}

func TestRedefinesBuiltinIDAfterGo1_21(t *testing.T) {
	testRule(t, "go1.21/redefines_builtin_id", &rule.RedefinesBuiltinIDRule{})
}
