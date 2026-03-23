package no_net_dial_timeout_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_net_dial_timeout"
)

func TestNoNetDialTimeout(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_dial_timeout", &no_net_dial_timeout.NoNetDialTimeoutRule{})
}
