package no_ssrf_via_variable_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ssrf_via_variable"
)

func TestNoSsrfViaVariable(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ssrf_via_variable", &no_ssrf_via_variable.NoSsrfViaVariableRule{})
}
