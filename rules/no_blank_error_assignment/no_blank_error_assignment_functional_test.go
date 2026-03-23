package no_blank_error_assignment_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_blank_error_assignment"
)

func TestNoBlankErrorAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_blank_error_assignment", &no_blank_error_assignment.NoBlankErrorAssignmentRule{})
}
