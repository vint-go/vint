package use_switch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_switch"
)

func TestUseSwitch(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_switch", &use_switch.UseSwitchRule{})
}
