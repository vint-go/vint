package no_net_dial_timeout_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_dial_timeout"
)

func TestNoNetDialTimeout(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_dial_timeout", &no_net_dial_timeout.NoNetDialTimeoutRule{})
}
