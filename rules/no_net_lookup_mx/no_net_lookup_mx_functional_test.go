package no_net_lookup_mx_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_lookup_mx"
)

func TestNoNetLookupMx(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_mx", &no_net_lookup_mx.NoNetLookupMxRule{})
}
