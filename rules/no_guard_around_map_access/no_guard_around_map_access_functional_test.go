package no_guard_around_map_access_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_guard_around_map_access"
)

func TestNoGuardAroundMapAccess(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_guard_around_map_access", &no_guard_around_map_access.NoGuardAroundMapAccessRule{})
}
