package no_blank_error_assignment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_blank_error_assignment"
)

func TestNoBlankErrorAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_blank_error_assignment", &no_blank_error_assignment.NoBlankErrorAssignmentRule{})
}
