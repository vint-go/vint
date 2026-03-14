package no_smtp_injection_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_smtp_injection_taint"
)

func TestNoSmtpInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_smtp_injection_taint", &no_smtp_injection_taint.NoSmtpInjectionTaintRule{})
}
