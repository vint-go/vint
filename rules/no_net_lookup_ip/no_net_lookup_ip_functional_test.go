package no_net_lookup_ip_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_ip"
)

func TestNoNetLookupIp(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_ip", &no_net_lookup_ip.NoNetLookupIpRule{})
}
