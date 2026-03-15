package no_net_listen_packet_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_listen_packet"
)

func TestNoNetListenPacket(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_listen_packet", &no_net_listen_packet.NoNetListenPacketRule{})
}
