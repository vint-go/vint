package use_switch_default_style_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_switch_default_style"
)

func TestUseSwitchDefaultStyle(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_switch_default_style", &use_switch_default_style.EnforceSwitchStyleRule{})
}

func TestUseSwitchDefaultStyleAllowNoDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_switch_default_style_allow_no_default", &use_switch_default_style.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-no-default"},
	})
}

func TestUseSwitchDefaultStyleAllowNotLast(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_switch_default_style_allow_not_last", &use_switch_default_style.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-default-not-last"},
	})
}

func TestUseSwitchDefaultStyleAllowNoDefaultAllowNotLast(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_switch_default_style_allow_no_default_allow_not_last", &use_switch_default_style.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-no-default", "allow-default-not-last"},
	})
}
