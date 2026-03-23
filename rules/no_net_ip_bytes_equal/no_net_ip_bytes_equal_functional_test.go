package no_net_ip_bytes_equal_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_ip_bytes_equal"
)

func TestNoNetIpBytesEqual(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_ip_bytes_equal", &no_net_ip_bytes_equal.NoNetIpBytesEqualRule{})
}
