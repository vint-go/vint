package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestUseWaitGroupGo(t *testing.T) {
	testRule(t, "use_waitgroup_go", &rule.UseWaitGroupGoRule{}, &lint.RuleConfig{})
	testRule(t, "go1.25/use_waitgroup_go", &rule.UseWaitGroupGoRule{}, &lint.RuleConfig{})
}
