package no_tls_session_resumption_bypass_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_tls_session_resumption_bypass"
)

func TestNoTlsSessionResumptionBypass(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_tls_session_resumption_bypass", &no_tls_session_resumption_bypass.NoTlsSessionResumptionBypassRule{})
}
