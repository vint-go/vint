package no_command_injection_taint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_command_injection_taint"
)

func TestNoCommandInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_command_injection_taint", &no_command_injection_taint.NoCommandInjectionTaintRule{})
}
