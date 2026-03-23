package no_suspicious_map_key_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_suspicious_map_key"
)

func TestNoSuspiciousMapKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_suspicious_map_key", &no_suspicious_map_key.NoSuspiciousMapKeyRule{})
}
