package use_type_switch_guard_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_type_switch_guard"
)

func TestUseTypeSwitchGuard(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_switch_guard", &use_type_switch_guard.UseTypeSwitchGuardRule{})
}
