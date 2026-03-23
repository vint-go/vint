package use_type_switch_chain_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_type_switch_chain"
)

func TestUseTypeSwitchChain(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_switch_chain", &use_type_switch_chain.UseTypeSwitchChainRule{})
}
