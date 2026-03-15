package no_net_lookup_addr_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_addr"
)

func TestNoNetLookupAddr(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_addr", &no_net_lookup_addr.NoNetLookupAddrRule{})
}
