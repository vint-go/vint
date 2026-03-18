package no_weak_slice_guard_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_weak_slice_guard"
)

func TestNoWeakSliceGuard(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_weak_slice_guard", &no_weak_slice_guard.NoWeakSliceGuardRule{})
}
