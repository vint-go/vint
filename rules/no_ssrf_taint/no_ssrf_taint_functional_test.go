package no_ssrf_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_ssrf_taint"
)

func TestNoSsrfTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ssrf_taint", &no_ssrf_taint.NoSsrfTaintRule{})
}
