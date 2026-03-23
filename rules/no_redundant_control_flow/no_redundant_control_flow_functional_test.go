package no_redundant_control_flow_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_control_flow"
)

func TestNoRedundantControlFlow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_control_flow", &no_redundant_control_flow.NoRedundantControlFlowRule{})
}
