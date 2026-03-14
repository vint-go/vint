package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestForbiddenCallInWgGo(t *testing.T) {
	testRule(t, "forbidden_call_in_wg_go", &rule.ForbiddenCallInWgGoRule{}, &lint.RuleConfig{})
	testRule(t, "go1.25/forbidden_call_in_wg_go", &rule.ForbiddenCallInWgGoRule{}, &lint.RuleConfig{})
}
