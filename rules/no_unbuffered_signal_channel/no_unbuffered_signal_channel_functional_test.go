package no_unbuffered_signal_channel_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unbuffered_signal_channel"
)

func TestNoUnbufferedSignalChannel(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unbuffered_signal_channel", &no_unbuffered_signal_channel.NoUnbufferedSignalChannelRule{})
}
