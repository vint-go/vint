package no_ssrf_taint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ssrf_taint"
)

func TestNoSsrfTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ssrf_taint", &no_ssrf_taint.NoSsrfTaintRule{})
}
