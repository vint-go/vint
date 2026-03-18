package no_redundant_control_flow_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_control_flow"
)

func TestNoRedundantControlFlow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_control_flow", &no_redundant_control_flow.NoRedundantControlFlowRule{})
}
