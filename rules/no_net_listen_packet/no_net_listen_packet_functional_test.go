package no_net_listen_packet_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_listen_packet"
)

func TestNoNetListenPacket(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_listen_packet", &no_net_listen_packet.NoNetListenPacketRule{})
}
