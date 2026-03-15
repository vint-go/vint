package no_mismatched_append_assign_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_mismatched_append_assign"
)

func TestNoMismatchedAppendAssign(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_mismatched_append_assign", &no_mismatched_append_assign.NoMismatchedAppendAssignRule{})
}
