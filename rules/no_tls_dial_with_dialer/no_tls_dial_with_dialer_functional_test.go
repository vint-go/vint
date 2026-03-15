package no_tls_dial_with_dialer_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_tls_dial_with_dialer"
)

func TestNoTlsDialWithDialer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_tls_dial_with_dialer", &no_tls_dial_with_dialer.NoTlsDialWithDialerRule{})
}
