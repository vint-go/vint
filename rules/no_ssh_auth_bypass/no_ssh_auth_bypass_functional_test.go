package no_ssh_auth_bypass_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_ssh_auth_bypass"
)

func TestNoSshAuthBypass(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ssh_auth_bypass", &no_ssh_auth_bypass.NoSshAuthBypassRule{})
}
