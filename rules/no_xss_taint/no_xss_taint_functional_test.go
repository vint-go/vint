package no_xss_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_xss_taint"
)

func TestNoXssTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_xss_taint", &no_xss_taint.NoXssTaintRule{})
}
