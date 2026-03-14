package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestUncheckedDynamicCast(t *testing.T) {
	testRule(t, "unchecked_type_assertion", &rule.UncheckedTypeAssertionRule{})
}

func TestUncheckedDynamicCastWithAcceptIgnored(t *testing.T) {
	testRule(t, "unchecked_type_assertion_accept_ignored", &rule.UncheckedTypeAssertionRule{},
		&lint.RuleConfig{
			Arguments: lint.Arguments{
				map[string]any{"acceptIgnoredAssertionResult": true},
			},
		},
	)
	testRule(t, "unchecked_type_assertion_accept_ignored", &rule.UncheckedTypeAssertionRule{},
		&lint.RuleConfig{
			Arguments: lint.Arguments{
				map[string]any{"accept-ignored-assertion-result": true},
			},
		},
	)
}
