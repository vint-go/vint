package no_self_assignment_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_self_assignment"
)

func TestNoSelfAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_self_assignment", &no_self_assignment.NoSelfAssignmentRule{})
}
