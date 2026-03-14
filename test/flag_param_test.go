package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestFlagParam(t *testing.T) {
	testRule(t, "flag_param", &rule.FlagParamRule{})
}
