package no_magic_number_in_assignment_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_magic_number_in_assignment"
)

func TestNoMagicNumberInAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_assignment", &no_magic_number_in_assignment.NoMagicNumberInAssignmentRule{})
}
