package use_tagged_switch_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_tagged_switch"
)

func TestUseTaggedSwitch(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_tagged_switch", &use_tagged_switch.UseTaggedSwitchRule{})
}
