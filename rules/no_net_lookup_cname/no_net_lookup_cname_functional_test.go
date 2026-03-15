package no_net_lookup_cname_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_cname"
)

func TestNoNetLookupCname(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_cname", &no_net_lookup_cname.NoNetLookupCnameRule{})
}
