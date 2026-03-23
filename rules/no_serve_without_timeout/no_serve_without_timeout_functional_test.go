package no_serve_without_timeout_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_serve_without_timeout"
)

func TestNoServeWithoutTimeout(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_serve_without_timeout", &no_serve_without_timeout.NoServeWithoutTimeoutRule{})
}
