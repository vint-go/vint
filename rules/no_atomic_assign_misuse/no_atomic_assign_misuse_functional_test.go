package no_atomic_assign_misuse_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_atomic_assign_misuse"
)

func TestNoAtomicAssignMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_atomic_assign_misuse", &no_atomic_assign_misuse.NoAtomicAssignMisuseRule{})
}
