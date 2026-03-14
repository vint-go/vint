package no_log_injection_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_log_injection_taint"
)

func TestNoLogInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_log_injection_taint", &no_log_injection_taint.NoLogInjectionTaintRule{})
}
