package no_invalid_host_port_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_host_port"
)

func TestNoInvalidHostPort(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_host_port", &no_invalid_host_port.NoInvalidHostPortRule{})
}
