package no_command_injection_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_command_injection_taint"
)

func TestNoCommandInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_command_injection_taint", &no_command_injection_taint.NoCommandInjectionTaintRule{})
}
