package no_weak_slice_guard_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_weak_slice_guard"
)

func TestNoWeakSliceGuard(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_weak_slice_guard", &no_weak_slice_guard.NoWeakSliceGuardRule{})
}
