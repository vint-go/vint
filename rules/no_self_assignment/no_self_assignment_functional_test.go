package no_self_assignment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_self_assignment"
)

func TestNoSelfAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_self_assignment", &no_self_assignment.NoSelfAssignmentRule{})
}
