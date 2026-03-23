package use_loop_condition_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_loop_condition"
)

func TestUseLoopCondition(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_loop_condition", &use_loop_condition.UseLoopConditionRule{})
}
