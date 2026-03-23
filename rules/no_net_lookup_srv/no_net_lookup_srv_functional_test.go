package no_net_lookup_srv_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_lookup_srv"
)

func TestNoNetLookupSrv(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_srv", &no_net_lookup_srv.NoNetLookupSrvRule{})
}
