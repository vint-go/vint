package no_slog_key_value_mismatch_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_slog_key_value_mismatch"
)

func TestNoSlogKeyValueMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_slog_key_value_mismatch", &no_slog_key_value_mismatch.NoSlogKeyValueMismatchRule{})
}
