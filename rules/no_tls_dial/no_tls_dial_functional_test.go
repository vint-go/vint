package no_tls_dial_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_tls_dial"
)

func TestNoTlsDial(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_tls_dial", &no_tls_dial.NoTlsDialRule{})
}
