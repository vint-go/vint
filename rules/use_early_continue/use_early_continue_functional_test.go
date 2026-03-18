package use_early_continue_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_early_continue"
)

func TestUseEarlyContinue(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_early_continue", &use_early_continue.UseEarlyContinueRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}
