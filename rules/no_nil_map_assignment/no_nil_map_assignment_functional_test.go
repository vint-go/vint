package no_nil_map_assignment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nil_map_assignment"
)

func TestNoNilMapAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nil_map_assignment", &no_nil_map_assignment.NoNilMapAssignmentRule{})
}
