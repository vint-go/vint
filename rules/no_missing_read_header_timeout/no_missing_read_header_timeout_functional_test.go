package no_missing_read_header_timeout_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_missing_read_header_timeout"
)

func TestNoMissingReadHeaderTimeout(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_missing_read_header_timeout", &no_missing_read_header_timeout.NoMissingReadHeaderTimeoutRule{})
}
