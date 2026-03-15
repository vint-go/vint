package no_net_lookup_port_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_port"
)

func TestNoNetLookupPort(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_port", &no_net_lookup_port.NoNetLookupPortRule{})
}
