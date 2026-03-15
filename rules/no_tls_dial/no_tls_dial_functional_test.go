package no_tls_dial_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_tls_dial"
)

func TestNoTlsDial(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_tls_dial", &no_tls_dial.NoTlsDialRule{})
}
