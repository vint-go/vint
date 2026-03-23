package no_ineffectual_assignment_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ineffectual_assignment"
)

func TestNoIneffectualAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ineffectual_assignment", &no_ineffectual_assignment.NoIneffectualAssignmentRule{})
}
