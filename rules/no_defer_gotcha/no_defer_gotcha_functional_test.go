package no_defer_gotcha_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_defer_gotcha"
)

func TestNoDeferGotcha(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_gotcha", &no_defer_gotcha.DeferRule{})
}

func TestNoDeferGotchaLoopDisabled(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_gotcha_loop_disabled", &no_defer_gotcha.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"return", "recover", "callChain", "methodCall"}},
	})
	functional_test_helpers.TestRule(t, "no_defer_gotcha_loop_disabled", &no_defer_gotcha.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"return", "recover", "call-chain", "method-call"}},
	})
}

func TestNoDeferGotchaOthersDisabled(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_gotcha_only_loop_enabled", &no_defer_gotcha.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"loop"}},
	})
}
