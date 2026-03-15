package no_atomic_assign_misuse_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_atomic_assign_misuse"
)

func TestNoAtomicAssignMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_atomic_assign_misuse", &no_atomic_assign_misuse.NoAtomicAssignMisuseRule{})
}
