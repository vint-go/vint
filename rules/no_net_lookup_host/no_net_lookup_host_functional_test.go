package no_net_lookup_host_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_host"
)

func TestNoNetLookupHost(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_host", &no_net_lookup_host.NoNetLookupHostRule{})
}
