package no_net_lookup_ns_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_lookup_ns"
)

func TestNoNetLookupNs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_ns", &no_net_lookup_ns.NoNetLookupNsRule{})
}
