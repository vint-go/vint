package use_early_return_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_early_return"
)

func TestUseEarlyReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_early_return", &use_early_return.EarlyReturnRule{})
	functional_test_helpers.TestRule(t, "use_early_return_scope", &use_early_return.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	functional_test_helpers.TestRule(t, "use_early_return_scope", &use_early_return.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
	functional_test_helpers.TestRule(t, "use_early_return_jump", &use_early_return.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allowJump"}})
	functional_test_helpers.TestRule(t, "use_early_return_jump", &use_early_return.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allow-jump"}})
	functional_test_helpers.TestRule(t, "use_early_return_jump_scope", &use_early_return.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allow-jump", "preserve-scope"}})
}
