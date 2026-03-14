package no_short_rsa_key_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_short_rsa_key"
)

func TestNoShortRsaKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_short_rsa_key", &no_short_rsa_key.NoShortRsaKeyRule{})
}
