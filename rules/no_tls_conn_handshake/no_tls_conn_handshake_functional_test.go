package no_tls_conn_handshake_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_tls_conn_handshake"
)

func TestNoTlsConnHandshake(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_tls_conn_handshake", &no_tls_conn_handshake.NoTlsConnHandshakeRule{})
}
