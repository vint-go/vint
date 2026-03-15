package no_net_dial_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_dial"
)

func TestNoNetDial(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_dial", &no_net_dial.NoNetDialRule{})
}
