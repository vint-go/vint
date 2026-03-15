package no_unsafe_cors_bypass_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unsafe_cors_bypass"
)

func TestNoUnsafeCorsBypass(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsafe_cors_bypass", &no_unsafe_cors_bypass.NoUnsafeCorsBypassRule{})
}
