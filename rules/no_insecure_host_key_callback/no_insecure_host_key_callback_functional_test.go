package no_insecure_host_key_callback_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_insecure_host_key_callback"
)

func TestNoInsecureHostKeyCallback(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_insecure_host_key_callback", &no_insecure_host_key_callback.NoInsecureHostKeyCallbackRule{})
}
