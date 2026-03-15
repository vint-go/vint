package no_net_lookup_ns_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_ns"
)

func TestNoNetLookupNs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_ns", &no_net_lookup_ns.NoNetLookupNsRule{})
}
