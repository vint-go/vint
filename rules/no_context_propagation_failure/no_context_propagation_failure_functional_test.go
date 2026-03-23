package no_context_propagation_failure_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_context_propagation_failure"
)

func TestNoContextPropagationFailure(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_context_propagation_failure", &no_context_propagation_failure.NoContextPropagationFailureRule{})
}
