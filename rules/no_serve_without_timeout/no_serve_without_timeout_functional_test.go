package no_serve_without_timeout_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_serve_without_timeout"
)

func TestNoServeWithoutTimeout(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_serve_without_timeout", &no_serve_without_timeout.NoServeWithoutTimeoutRule{})
}
